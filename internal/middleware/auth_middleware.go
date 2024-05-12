package middleware

//"github.com/Viet-ph/Furniture-Store-Server/internal/service"

type ContextKey string

const ContextUserKey ContextKey = "user"

// func NewMiddlewareAuth(userService *service.UserService) func(http.Handler) http.Handler {
// 	//This will take the dependencies and return a authentication middleware that accepts only a single handler.
// 	//By doing this, will clean up the middleware function arguments and create closure to outter deps.
// 	return func(handler http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			tokenString := ExtractTokenFromHeader(r)
// 			token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, errors.New("unexpected signing method")
// 				}
// 				return []byte(os.Getenv("JWT_SECRET")), nil
// 			})
// 			if err != nil || !token.Valid {
// 				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized Access")
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), "token", token)
// 			handler.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }
