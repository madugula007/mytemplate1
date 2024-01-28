package supertest

func (ctx *Supertest) Set(key, value string) {
	ctx.request.httpRequest.Header.Set(key, value)
}
