package committer

import (
	"context"

	"cloud.google.com/go/spanner"
)

type CommitPlan struct {
	mutations []*spanner.Mutation
}

func NewPlan() *CommitPlan {
	return &CommitPlan{mutations: make([]*spanner.Mutation, 0)}
}

func (p *CommitPlan) Add(m *spanner.Mutation) {
	if m != nil {
		p.mutations = append(p.mutations, m)
	}
}

func (p *CommitPlan) Mutations() []*spanner.Mutation {
	return p.mutations
}

func (p *CommitPlan) IsEmpty() bool {
	return len(p.mutations) == 0
}

type Committer struct {
	client *spanner.Client
}

func New(client *spanner.Client) *Committer {
	return &Committer{client: client}
}

func (c *Committer) Apply(ctx context.Context, plan *CommitPlan) error {
	if plan.IsEmpty() {
		return nil
	}

	_, err := c.client.Apply(ctx, plan.Mutations())
	return err
}
