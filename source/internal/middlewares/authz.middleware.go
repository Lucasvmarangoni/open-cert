package middlewares

// import (
// 	"net/http"
// 	"slices"

// 	p "github.com/Lucasvmarangoni/financial-file-manager/config/casbin"
// 	"github.com/casbin/casbin/v2"
// 	"github.com/go-chi/jwtauth"
// )

// type Authorization struct {
// 	enforcer *casbin.Enforcer
// 	policy   string
// 	model    string
// 	admin    []string
// 	read     []string
// 	rules    [][]string
// }

// func NewAuthorization(policy string, model string) *Authorization {
// 	a := &Authorization{
// 		policy: policy,
// 		model:  model,
// 	}
// 	a.init()
// 	return a
// }

// func (a *Authorization) init() {
// 	policy := p.NewPolice()
// 	policy.SetPolicy()

// 	a.rules = policy.Rules
// 	a.admin = policy.Admin
// 	a.read = policy.Read

// 	a.enforcer, _ = casbin.NewEnforcer(a.model)

// 	for _, rule := range a.rules {
// 		a.enforcer.AddPolicy(rule)
// 	}

// 	for _, rule := range policy.Groups {
// 		a.enforcer.AddGroupingPolicy(rule)
// 	}

// }

// func (a *Authorization) Authorizer() func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			var role string
// 			id := getSub(w, r)

// 			if slices.Contains(a.admin, id) || slices.Contains(a.read, id) {
// 				role = id
// 			} else {
// 				role = "member"
// 			}
// 			method := r.Method
// 			path := r.URL.Path
		
// 			if ok, _ := a.enforcer.Enforce(role, path, method); !ok {
// 				http.Error(w, "Forbidden", http.StatusForbidden)
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

