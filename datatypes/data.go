package datatypes

type Data interface {
  String() string
  Type() string
  Response() []byte
}
