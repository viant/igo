package metric

type Stmt uint64

const (
	stmtNone    = Stmt(0)
	stmtAssign  = Stmt(1)
	stmtDeclare = Stmt(1 << 2)
	stmtExpr    = Stmt(1 << 3)
	stmtIfElse  = Stmt(1 << 4)
	stmtFor     = Stmt(1 << 5)
	stmtRange   = Stmt(1 << 6)
	stmtCall    = Stmt(1 << 7)
)

func (s *Stmt) Reset() {
	*s = stmtNone
}

func (s *Stmt) FlagDeclare() {
	*s |= stmtDeclare
}

func (s *Stmt) FlagAssign() {
	*s |= stmtAssign
}

func (s *Stmt) FlagExpr() {
	*s |= stmtExpr
}

func (s *Stmt) FlagIfElse() {
	*s |= stmtIfElse
}

func (s *Stmt) FlagFor() {
	*s |= stmtFor
}

func (s *Stmt) FlagRange() {
	*s |= stmtRange
}

func (s *Stmt) FlagCall() {
	*s |= stmtCall
}

func (s Stmt) UsesDeclare() bool {
	return s&stmtDeclare != 1
}

func (s Stmt) UsesAssign() bool {
	return s&stmtAssign != 1
}

func (s Stmt) UsesExpr() bool {
	return s&stmtExpr != 1
}

func (s Stmt) UsesIfElse() bool {
	return s&stmtIfElse != 1
}

func (s Stmt) UsesFor() bool {
	return s&stmtFor != 1
}

func (s Stmt) UsesRange() bool {
	return s&stmtRange != 1
}

func (s Stmt) UsesCall() bool {
	return s&stmtCall != 1
}
