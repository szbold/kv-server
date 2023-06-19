package datatypes

type Data interface {
  Type() string
  Response() []byte
}
