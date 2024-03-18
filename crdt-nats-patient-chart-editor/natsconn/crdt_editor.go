package natsconn

import "errors"

//들어오는 메세지를 받아서, 하나의 차트 object로 merge할 수 있어야 한다.
//그리고 그걸 다시 클라이언트에게 publish해야 한다.

type Chart struct{}

func CRDTEditChart(chartEventData []byte) (*Chart, error) {
	return nil, errors.New("not implemented")
}
