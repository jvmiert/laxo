package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/shop"
)

type testHandlerService struct {
	lazada *lazada.Service
	shop   *shop.Service
	assets *assets.Service
}

type testHandler struct {
	server  *laxo.Server
	service *testHandlerService
}

func InitTestHandler(server *laxo.Server, l *lazada.Service, p *shop.Service, a *assets.Service, r *mux.Router, n *negroni.Negroni) {
	s := &testHandlerService{
		lazada: l,
		shop:   p,
		assets: a,
	}

	h := testHandler{
		server:  server,
		service: s,
	}

	r.Handle("/test/test", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleTest)),
	)).Methods("GET")
}

func (h *testHandler) HandleTest(w http.ResponseWriter, r *http.Request, uID string) {
	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = h.service.shop.GetProductDetailsByID("01G5TV4071TZES7SVTTB1908WR", s.Model.ID)
	if errors.Is(err, shop.ErrProductNotFound) {
		laxo.ErrorJSONEncode(w, err, http.StatusNotFound)
		return
	}

	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

  test := `<div><h2>Giới thiệu về Android Tivi Sony 4K 49 inch KD-49X8500G/S</h2> <h3>Thiết kế hài hòa, chất lượng</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S</strong> mẫu được sản xuất năm 2019 thiết kế với kiểu dáng tinh tế, đẹp mắt. Viền màn hình mỏng được làm từ kim loại sang trọng tăng thêm vẻ đẹp cho không gian bạn sử dụng. Chân đế được cấu tạo kim loại chắc chắn, đẹp mắt giúp tivi có thể đứng vững trên mọi mặt phẳng .Kích thước màn hình tivi lớn 49 inch phù hợp với những không gian phòng rộng như: phòng khách, phòng họp công ty, sảnh khách sạn, lớp học,…</p> <p><img alt="Android Tivi Sony 4K 49 inch KD-49X8500G/S phù hợp với không gian trong phòng khách" height="444" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/dbcc6e36e01fde4bc7405b731a423e81.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/06/chuẩn-1-1024x614.jpg 1024w, https://tmp.phongvu.vn/wp-content/uploads/2019/06/chuẩn-1-300x180.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/06/chuẩn-1-768x461.jpg 768w, https://tmp.phongvu.vn/wp-content/uploads/2019/06/chuẩn-1.jpg 1200w" width="740"></p> <h3>Độ phân giải 4K đem tới hình ảnh sắc nét</h3> <p><strong>Android Tivi Sony KD-49X8500G/S </strong>với độ phân giải 4K gấp 4 lần Full HD mang tới hình ảnh vô cùng sắc nét, chân thật cho người xem.</p> <p><img alt="Android Tivi Sony 4K 49 inch KD-49X8500G/S độ phân giải 4k đem tới hình ảnh sắc nét chân thực nhất" height="411" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/d93baaa57f0034f686166da200267b1d.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/06/4k-789.jpg 780w, https://tmp.phongvu.vn/wp-content/uploads/2019/06/4k-789-300x167.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/06/4k-789-768x426.jpg 768w" width="740"></p> <h3>Công nghệ hình ảnh HDR10 đặc sắc</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S </strong>được tích hợp công nghệ HDR10 làm tăng thêm độ tương phản và màu sắc cho hình ảnh giúp tivi hiển thị hình ảnh ở các vùng tối, sáng được rõ nét mang đến độ sâu, vô cùng sắc nét giúp người xem trải nghiệm hình ảnh tốt nhất</p> <p><img alt="Hdr10 giúp tái tạo chất lượng hình ảnh đem hình ảnh chi tiết sấc nét" height="433" sizes="(max-width: 780px) 100vw, 780px" src="https://vn-live-02.slatic.net/p/324e30c4d556d7006715fe898a6c10d0.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/Hdr10-1.jpg 780w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/Hdr10-1-300x167.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/Hdr10-1-768x426.jpg 768w" width="780"></p> <h3>Công nghệ đèn nền Edge LED Frame Dimming</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S </strong>đem tới chiều sâu cho hình ảnh tăng trải nghiệm cho người sử dụng.</p> <p><img alt="Edge LED Frame Dimming đem tới hình ảnh có độ sáng chất lượng tuyệt vời" height="411" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/0e0032b2e407cba5883db7fa001ba8c0.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/Edge-LED-Frame-Dimming.jpg 780w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/Edge-LED-Frame-Dimming-300x167.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/Edge-LED-Frame-Dimming-768x426.jpg 768w" width="740"></p> <h3>Công nghệ hình ảnh 4K X-Reality Pro</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S</strong> được ứng dụng công nghệ độc quyền của hãng giúp nâng cấp chất lượng hình ảnh với độ sắc nét cao, mượt mà, hạn chế hình ảnh bị nhòe, nhiễu hạt</p> <p><img alt="X-Reality Pro độc quyền của hãng giúp nâng cấp chất lượng hình ảnh với độ sắc nét cao, mượt mà, hạn chế hình ảnh bị nhòe, nhiễu hạt" height="372" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/4a5c303d4fd0d082ca2a74f616a3cd46.png" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/X-Reality-Pro-2-1.png 730w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/X-Reality-Pro-2-1-300x151.png 300w" width="740"></p> <h3>Công nghệ hình ảnh TRILUMINOS đặc biệt</h3> <p>Công nghệ TRILUMINOS giúp tái tạo những mảng màu khó tái tạo nhất trở nên rực rỡ, tự nhiên. Những màu xanh lục, xanh lam và đỏ được thể hiện rực rỡ đem lại trải nghiệm hình ảnh giống ngoài tự nhiện</p> <p><img alt="công nghệ TRILUMINOS đem tới dải màu rộng sống động chân thực" height="339" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/d47fabdd8e6b94965fc97453beeaf1c0.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/01/triluminos-2.jpg 730w, https://tmp.phongvu.vn/wp-content/uploads/2019/01/triluminos-2-300x137.jpg 300w" width="740"></p> <h3>Bộ xử lý 4K HDR Processor X1</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S </strong>sở hữu bộ xử lý 4K HDR Processor X1 tái tạo chiều sâu, kết cấu và màu tự nhiên ở mức độ nâng cao cho chất lượng hình ảnh trung thực tuyệt đối.</p> <p><img alt="bộ xử lý 4K HDR Processor X1 tái tạo chiều sâu, kết cấu và màu tự nhiên" height="381" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/bd5902d15aef1bb2aa94f6114715397a.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-42-e1557285496458.jpg 1428w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-42-e1557285496458-300x154.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-42-e1557285496458-768x395.jpg 768w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-42-e1557285496458-1024x527.jpg 1024w" width="740"></p> <h3>Âm thanh mạnh mẽ, tuyệt vời</h3> <p>Tivi được thiết kế loa tiên tiến 20W cùng công nghệ âm thanh tiên tiến nhất Acoustic Multi-Audio, hỗ trợ eARC, S-Force Front Surround đem tớ chất lượng âm thanh vòm kỹ thuật số theo nhiều kênh riêng biệt đem tới âm thanh đa chiều với độ sâu vang tạo cảm giác như đang ở rạp hát, rạp chiếu phim.</p> <p><img alt="Android Tivi Sony 4K 49 inch KD-49X8500G/S với âm thanh sống động mạnh mẽ" height="411" src="https://vn-live-02.slatic.net/p/bee977a1dec8aef75b4ca8cb3e3a2ebf.jpg" width="740"></p> <h3>Hệ điều hành Android phong phú, dễ sử dụng mọi đối tượng</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S</strong> chạy hệ điều hành Android với giao diện đẹp mắt thân thiện dễ dàng cho người dùng. Tivi được trang bị sẵn các ứng dụng vô cùng tiện ích: Trình duyệt Web, YouTube, Netflix,…Cùng với đó người dùng cài đặt thêm những ứng dụng khác: Fim+, Asphalt 8, ZingMp3, NCT,…</p> <p><img alt="Android Tivi Sony 4K 49 inch KD-49X8500G/S với hệ điều hành Android đem tới chất lượng cho người sử dụng" height="411" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/92032b7e7f6553ac98bbe6555c9b140b.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-43.jpg 780w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-43-300x167.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/1-43-768x426.jpg 768w" width="740"></p> <h3>Hỗ trợ tìm kiếm bằng giọng nói</h3> <p><strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S </strong>được tích hợp tính năng cực kỳ hữu ích giúp người dùng dễ dàng tìm kiếm các ứng dụng, tên các bộ phim cực kỳ nhanh chóng và hiệu quả hơn cùng gia đình và bạn bè cùng tận hưởng và giải trí sau những ngày làm việc và học tập căng thẳng.</p> <p><img alt="Android Tivi Sony 4K 49 inch KD-49X8500G/S hỗ trợ tìm kiếm bằng giọng nói" height="444" sizes="(max-width: 740px) 100vw, 740px" src="https://vn-live-02.slatic.net/p/bda2943834a0c04a298f22d3a48fa3aa.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/tivi-4k-sony-kd-55x8500g-5.jpg 677w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/tivi-4k-sony-kd-55x8500g-5-300x180.jpg 300w" width="740"></p> <h3>Đa dạng cổng kết nối <strong>Android Tivi Sony 4K 49 inch KD-49X8500G/S</strong></h3> <p>Tivi được thiết kế nhiều cổng kết nối tiện cho người sử dụng kết nối các thiết bị bên ngoài nhanh chóng và thuận tiện nhất.</p> <p>– HDMI: Kết nối tivi với laptop, pc, dàn âm thanh để trình chiếu các nội dung hình ảnh, ca nhạc.</p> <p>– USB: Kết nối và phát trực tiếp nội dung có sẵn trên USB.</p> <p>– Kết nối mạng LAN, Wifi để kết nối Internet giúp bạn dễ dàng xem xem Youtube, truy cập Facebook, đọc báo</p> <p>– Cổng Optical, HDMI ARC để xuất âm thanh từ tivi xuống loa.</p> <p><img alt="Android Tivi Sony 4K 49 inch KD-49X8500G/S với cổng kết nối đa dạng" height="433" sizes="(max-width: 780px) 100vw, 780px" src="https://vn-live-02.slatic.net/p/f4fccf789ed4bbd7dc564902e60280fc.jpg" srcset="https://tmp.phongvu.vn/wp-content/uploads/2019/05/thiết-sàd.jpg 780w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/thiết-sàd-300x167.jpg 300w, https://tmp.phongvu.vn/wp-content/uploads/2019/05/thiết-sàd-768x426.jpg 768w" width="780"></p></div>`

	findImages := regexp.MustCompile(`(http(s?):)([\/|.|\p{L}|\d|\s|-])*\.(?:jpe?g|gif|png)`)
	matches := findImages.FindAllStringSubmatch(test, -1)

	for _, element := range matches {
		h.server.Logger.Debugw("found image", "element", element[0])
	}

	_ = `<div style="margin:0"><span style="font-family:none"><strong style="font-weight:bold;font-family:none">Test</strong></span></div><div style="margin:0"><span style="font-family:none"></span></div><div style="margin:0"><span style="font-family:none"><em>Test</em></span></div><div style="margin:0"><span style="font-family:none"></span></div><div style="margin:0"><span style="font-family:none"><u>Test</u></span><i><span>Test</span></i><span> test</span></div><div style="margin:0"><span></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span><strong style="font-weight:bold">This is bold</strong></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span><u>This is underlined</u></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span><em>This is italic</em></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span></span></div><div style="margin:0"><h1><span>Heading1</span></h1></div><div style="margin:0"><h2><span>Heading2</span></h2></div><div style="margin:0"><h3><span>Heading3</span></h3></div><div style="margin:0"><span></span></div>`

	schema, _ := h.service.shop.HTMLToSlate(test)

	b, _ := json.Marshal(schema)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
