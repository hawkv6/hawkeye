package domain

type Sid interface {
	GetKey() string
	GetIgpRouterId() string
	GetSid() string
}
