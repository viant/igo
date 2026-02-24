package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
)

func (s *Scope) tryCompileMergeSortLikeForStmt(scope *Scope, forStmt *ast.ForStmt, init et.New, labelID uint32) (et.New, bool, error) {
	if labelID != 0 || forStmt.Cond == nil || forStmt.Post != nil || forStmt.Body == nil {
		return nil, false, nil
	}
	if p, ok := parseMergeCompareLoop(forStmt.Cond, forStmt.Body); ok {
		iSel, ok, err := selectDirectInt(scope, p.i)
		if !ok || err != nil {
			return nil, false, err
		}
		jSel, ok, err := selectDirectInt(scope, p.j)
		if !ok || err != nil {
			return nil, false, err
		}
		kSel, ok, err := selectDirectInt(scope, p.k)
		if !ok || err != nil {
			return nil, false, err
		}
		midSel, ok, err := selectDirectInt(scope, p.mid)
		if !ok || err != nil {
			return nil, false, err
		}
		rightSel, ok, err := selectDirectInt(scope, p.right)
		if !ok || err != nil {
			return nil, false, err
		}
		dataSel, ok, err := selectDirectIntSlice(scope, p.data)
		if !ok || err != nil {
			return nil, false, err
		}
		tmpSel, ok, err := selectDirectIntSlice(scope, p.tmp)
		if !ok || err != nil {
			return nil, false, err
		}
		newFor, err := et.NewForMergeCompareIntSliceStmt(init, iSel.Offset(), jSel.Offset(), kSel.Offset(), midSel.Offset(), rightSel.Offset(), dataSel.Offset(), tmpSel.Offset())
		return newFor, true, err
	}
	if p, ok := parseMergeCopyWithKLoop(forStmt.Cond, forStmt.Body); ok {
		idxSel, ok, err := selectDirectInt(scope, p.idx)
		if !ok || err != nil {
			return nil, false, err
		}
		limitSel, ok, err := selectDirectInt(scope, p.limit)
		if !ok || err != nil {
			return nil, false, err
		}
		kSel, ok, err := selectDirectInt(scope, p.k)
		if !ok || err != nil {
			return nil, false, err
		}
		dstSel, ok, err := selectDirectIntSlice(scope, p.dst)
		if !ok || err != nil {
			return nil, false, err
		}
		srcSel, ok, err := selectDirectIntSlice(scope, p.src)
		if !ok || err != nil {
			return nil, false, err
		}
		newFor, err := et.NewForMergeCopyWithKIntSliceStmt(init, idxSel.Offset(), limitSel.Offset(), kSel.Offset(), dstSel.Offset(), srcSel.Offset())
		return newFor, true, err
	}
	if p, ok := parseMergeCopyLoop(forStmt.Cond, forStmt.Body); ok {
		idxSel, ok, err := selectDirectInt(scope, p.idx)
		if !ok || err != nil {
			return nil, false, err
		}
		limitSel, ok, err := selectDirectInt(scope, p.limit)
		if !ok || err != nil {
			return nil, false, err
		}
		dstSel, ok, err := selectDirectIntSlice(scope, p.dst)
		if !ok || err != nil {
			return nil, false, err
		}
		srcSel, ok, err := selectDirectIntSlice(scope, p.src)
		if !ok || err != nil {
			return nil, false, err
		}
		newFor, err := et.NewForMergeCopyIntSliceStmt(init, idxSel.Offset(), limitSel.Offset(), dstSel.Offset(), srcSel.Offset())
		return newFor, true, err
	}
	return nil, false, nil
}

func selectDirectInt(scope *Scope, name string) (*exec.Selector, bool, error) {
	sel, err := scope.selector(&ast.Ident{Name: name}, false)
	if err != nil {
		return nil, false, nil
	}
	if sel.Pathway != exec.PathwayDirect || sel.Type.Kind() != reflect.Int {
		return nil, false, nil
	}
	return sel, true, nil
}

func selectDirectIntSlice(scope *Scope, name string) (*exec.Selector, bool, error) {
	sel, err := scope.selector(&ast.Ident{Name: name}, false)
	if err != nil {
		return nil, false, nil
	}
	if sel.Pathway != exec.PathwayDirect || sel.Type.Kind() != reflect.Slice || sel.Type.Elem().Kind() != reflect.Int {
		return nil, false, nil
	}
	return sel, true, nil
}

type mergeCompareLoopPattern struct {
	i     string
	j     string
	k     string
	mid   string
	right string
	data  string
	tmp   string
}

func parseMergeCompareLoop(cond ast.Expr, body *ast.BlockStmt) (mergeCompareLoopPattern, bool) {
	if body == nil || len(body.List) != 2 {
		return mergeCompareLoopPattern{}, false
	}
	andExpr, ok := unparenExpr(cond).(*ast.BinaryExpr)
	if !ok || andExpr.Op != token.LAND {
		return mergeCompareLoopPattern{}, false
	}
	i, mid, ok := parseLessIdentIdent(andExpr.X)
	if !ok {
		return mergeCompareLoopPattern{}, false
	}
	j, right, ok := parseLessIdentIdent(andExpr.Y)
	if !ok {
		return mergeCompareLoopPattern{}, false
	}
	ifStmt, ok := body.List[0].(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Else == nil || len(ifStmt.Body.List) != 2 {
		return mergeCompareLoopPattern{}, false
	}
	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok || len(elseBlock.List) != 2 {
		return mergeCompareLoopPattern{}, false
	}
	kInc, ok := body.List[1].(*ast.IncDecStmt)
	if !ok || kInc.Tok != token.INC {
		return mergeCompareLoopPattern{}, false
	}
	k, ok := identName(kInc.X)
	if !ok {
		return mergeCompareLoopPattern{}, false
	}
	cmp, ok := unparenExpr(ifStmt.Cond).(*ast.BinaryExpr)
	if !ok || cmp.Op != token.LEQ {
		return mergeCompareLoopPattern{}, false
	}
	leftSlice, leftIdx, ok := parseSliceIndex(cmp.X)
	if !ok || leftIdx != i {
		return mergeCompareLoopPattern{}, false
	}
	rightSlice, rightIdx, ok := parseSliceIndex(cmp.Y)
	if !ok || rightIdx != j || rightSlice != leftSlice {
		return mergeCompareLoopPattern{}, false
	}
	tmpSliceThen, tmpIdxThen, srcSliceThen, srcIdxThen, ok := parseAssignSliceToSlice(ifStmt.Body.List[0])
	if !ok || tmpIdxThen != k || srcSliceThen != leftSlice || srcIdxThen != i {
		return mergeCompareLoopPattern{}, false
	}
	thenInc, ok := ifStmt.Body.List[1].(*ast.IncDecStmt)
	if !ok || thenInc.Tok != token.INC {
		return mergeCompareLoopPattern{}, false
	}
	thenIncName, ok := identName(thenInc.X)
	if !ok || thenIncName != i {
		return mergeCompareLoopPattern{}, false
	}
	tmpSliceElse, tmpIdxElse, srcSliceElse, srcIdxElse, ok := parseAssignSliceToSlice(elseBlock.List[0])
	if !ok || tmpSliceElse != tmpSliceThen || tmpIdxElse != k || srcSliceElse != leftSlice || srcIdxElse != j {
		return mergeCompareLoopPattern{}, false
	}
	elseInc, ok := elseBlock.List[1].(*ast.IncDecStmt)
	if !ok || elseInc.Tok != token.INC {
		return mergeCompareLoopPattern{}, false
	}
	elseIncName, ok := identName(elseInc.X)
	if !ok || elseIncName != j {
		return mergeCompareLoopPattern{}, false
	}
	return mergeCompareLoopPattern{
		i:     i,
		j:     j,
		k:     k,
		mid:   mid,
		right: right,
		data:  leftSlice,
		tmp:   tmpSliceThen,
	}, true
}

type mergeCopyWithKPattern struct {
	idx   string
	limit string
	k     string
	dst   string
	src   string
}

func parseMergeCopyWithKLoop(cond ast.Expr, body *ast.BlockStmt) (mergeCopyWithKPattern, bool) {
	if body == nil || len(body.List) != 3 {
		return mergeCopyWithKPattern{}, false
	}
	idx, limit, ok := parseLessIdentIdent(cond)
	if !ok {
		return mergeCopyWithKPattern{}, false
	}
	dst, dstIdx, src, srcIdx, ok := parseAssignSliceToSlice(body.List[0])
	if !ok || srcIdx != idx {
		return mergeCopyWithKPattern{}, false
	}
	idxInc, ok := body.List[1].(*ast.IncDecStmt)
	if !ok || idxInc.Tok != token.INC {
		return mergeCopyWithKPattern{}, false
	}
	idxName, ok := identName(idxInc.X)
	if !ok || idxName != idx {
		return mergeCopyWithKPattern{}, false
	}
	kInc, ok := body.List[2].(*ast.IncDecStmt)
	if !ok || kInc.Tok != token.INC {
		return mergeCopyWithKPattern{}, false
	}
	kName, ok := identName(kInc.X)
	if !ok || kName != dstIdx {
		return mergeCopyWithKPattern{}, false
	}
	return mergeCopyWithKPattern{
		idx:   idx,
		limit: limit,
		k:     kName,
		dst:   dst,
		src:   src,
	}, true
}

type mergeCopyPattern struct {
	idx   string
	limit string
	dst   string
	src   string
}

func parseMergeCopyLoop(cond ast.Expr, body *ast.BlockStmt) (mergeCopyPattern, bool) {
	if body == nil || len(body.List) != 2 {
		return mergeCopyPattern{}, false
	}
	idx, limit, ok := parseLessIdentIdent(cond)
	if !ok {
		return mergeCopyPattern{}, false
	}
	dst, dstIdx, src, srcIdx, ok := parseAssignSliceToSlice(body.List[0])
	if !ok || dstIdx != idx || srcIdx != idx {
		return mergeCopyPattern{}, false
	}
	idxInc, ok := body.List[1].(*ast.IncDecStmt)
	if !ok || idxInc.Tok != token.INC {
		return mergeCopyPattern{}, false
	}
	idxName, ok := identName(idxInc.X)
	if !ok || idxName != idx {
		return mergeCopyPattern{}, false
	}
	return mergeCopyPattern{
		idx:   idx,
		limit: limit,
		dst:   dst,
		src:   src,
	}, true
}

func parseLessIdentIdent(expr ast.Expr) (string, string, bool) {
	binary, ok := unparenExpr(expr).(*ast.BinaryExpr)
	if !ok || binary.Op != token.LSS {
		return "", "", false
	}
	left, ok := identName(binary.X)
	if !ok {
		return "", "", false
	}
	right, ok := identName(binary.Y)
	if !ok {
		return "", "", false
	}
	return left, right, true
}

func identName(expr ast.Expr) (string, bool) {
	ident, ok := unparenExpr(expr).(*ast.Ident)
	if !ok {
		return "", false
	}
	return ident.Name, true
}

func parseAssignSliceToSlice(stmt ast.Stmt) (dstSlice, dstIdx, srcSlice, srcIdx string, ok bool) {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok || assign.Tok != token.ASSIGN || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return "", "", "", "", false
	}
	dstSlice, dstIdx, ok = parseSliceIndex(assign.Lhs[0])
	if !ok {
		return "", "", "", "", false
	}
	srcSlice, srcIdx, ok = parseSliceIndex(assign.Rhs[0])
	if !ok {
		return "", "", "", "", false
	}
	return dstSlice, dstIdx, srcSlice, srcIdx, true
}

func parseSliceIndex(expr ast.Expr) (string, string, bool) {
	indexExpr, ok := unparenExpr(expr).(*ast.IndexExpr)
	if !ok {
		return "", "", false
	}
	slice, ok := identName(indexExpr.X)
	if !ok {
		return "", "", false
	}
	index, ok := identName(indexExpr.Index)
	if !ok {
		return "", "", false
	}
	return slice, index, true
}
