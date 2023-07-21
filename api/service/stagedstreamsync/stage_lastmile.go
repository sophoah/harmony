package stagedstreamsync

import (
	"context"

	"github.com/harmony-one/harmony/core"
	"github.com/harmony-one/harmony/shard"
	"github.com/ledgerwatch/erigon-lib/kv"
)

type StageLastMile struct {
	configs StageLastMileCfg
}

type StageLastMileCfg struct {
	ctx context.Context
	bc  core.BlockChain
	db  kv.RwDB
}

func NewStageLastMile(cfg StageLastMileCfg) *StageLastMile {
	return &StageLastMile{
		configs: cfg,
	}
}

func NewStageLastMileCfg(ctx context.Context, bc core.BlockChain, db kv.RwDB) StageLastMileCfg {
	return StageLastMileCfg{
		ctx: ctx,
		bc:  bc,
		db:  db,
	}
}

func (lm *StageLastMile) Exec(ctx context.Context, firstCycle bool, invalidBlockRevert bool, s *StageState, reverter Reverter, tx kv.RwTx) (err error) {

	if s.state.initSync {
		return nil
	}

	// shouldn't execute for epoch chain
	if lm.configs.bc.ShardID() == shard.BeaconChainShardID && !s.state.isBeaconNode {
		return nil
	}

	bc := lm.configs.bc

	// update last mile blocks if any
	parentHash := bc.CurrentBlock().Hash()
	for {
		block := s.state.getBlockFromLastMileBlocksByParentHash(parentHash)
		if block == nil {
			break
		}
		err = s.state.UpdateBlockAndStatus(block, bc, false)
		if err != nil {
			break
		}
		parentHash = block.Hash()
	}

	return nil
}

func (lm *StageLastMile) Revert(ctx context.Context, firstCycle bool, u *RevertState, s *StageState, tx kv.RwTx) (err error) {
	useInternalTx := tx == nil
	if useInternalTx {
		tx, err = lm.configs.db.BeginRw(lm.configs.ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}

	if err = u.Done(tx); err != nil {
		return err
	}

	if useInternalTx {
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (lm *StageLastMile) CleanUp(ctx context.Context, firstCycle bool, p *CleanUpState, tx kv.RwTx) (err error) {
	useInternalTx := tx == nil
	if useInternalTx {
		tx, err = lm.configs.db.BeginRw(lm.configs.ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}

	if useInternalTx {
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}