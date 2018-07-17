package carrot

//Picker just pick message from queue without continuous consuming
type Picker struct {
	client *BrokerClient
}

//Pick item from queue
func (get *Picker) Pick(queue string) (*MessageContext, bool, error) {
	if ch, err := get.client.client.Channel(); err != nil {
		return nil, false, err
	} else {
		if msg, ok, err := ch.Get(queue, false); err != nil {
			return nil, false, err
		} else if ok {
			context := new(MessageContext)
			context.delivery = msg
			context.Message = Message{
				ContentType: msg.ContentType,
				Data:        msg.Body,
				Encoding:    msg.ContentEncoding,
				Headers:     msg.Headers,
			}
			return context, ok, nil
		} else {
			return nil, ok, nil
		}
	}
}

//NewPicker creates a new broker queue picker
func NewPicker(client *BrokerClient) *Picker {
	picker := new(Picker)
	picker.client = client
	return picker
}
