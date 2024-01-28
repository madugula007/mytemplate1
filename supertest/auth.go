package supertest

func (ctx *Supertest) Auth(username, password string) {
	ctx.request.httpRequest.SetBasicAuth(username, password)
}
