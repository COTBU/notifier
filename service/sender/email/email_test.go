package email

//
//func Test_sender_execTemplate(t *testing.T) {
//	type args struct {
//		data  any
//		templ string
//	}
//	tests := []struct {
//		name    string
//		sender  client
//		args    args
//		wantRes string
//		wantErr bool
//	}{
//		{
//			name:   "Test 1",
//			sender: senderMock,
//			args: args{
//				data: TemplateData{
//					Data: accountStatementDataNew,
//					IP:   "192.168.80.30",
//				},
//				templ: templateAccountStatementNew,
//			},
//			wantRes: result,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &client{
//				config:      tt.sender.config,
//				credentials: tt.sender.credentials,
//			}
//			gotRes, err := e.execTemplate(tt.args.data, tt.args.templ)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("execTemplate() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotRes != tt.wantRes {
//				t.Errorf("execTemplate() gotRes = %v, want %v", gotRes, tt.wantRes)
//			}
//		})
//	}
//}
