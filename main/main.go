package main

//func Auth(limiter *limiter.LoginLimiter) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		login := r.Header.Get("login")
//		bucket := limiter.GetBucket(login)
//		bucket.AddTokens()
//
//		if bucket.AvailableTokens() > 0 {
//			bucket.DeleteToken()
//			_, _ = w.Write([]byte("true"))
//			return
//		}
//		_, _ = w.Write([]byte("false"))
//	}
//}
//
//func main() {
//	limiter := limiter.NewLoginLimiter(10, time.Second*14)
//	http.HandleFunc("/auth", Auth(limiter))
//	err := http.ListenAndServe(":8888", nil)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
