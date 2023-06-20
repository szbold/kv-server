package types

type Data interface {
  Type() string
  Response() []byte
}
