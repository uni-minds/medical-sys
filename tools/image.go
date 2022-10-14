package tools

func CalcResize(w_ori, h_ori, w_target, h_target int) (w_resize, h_resize int) {
	ratio := float64(w_ori) / float64(h_ori)
	if w_target == 0 {
		h_resize = h_target
		w_resize = int(float64(h_target) * ratio)
	} else if h_target == 0 {
		w_resize = w_target
		h_resize = int(float64(w_target) / ratio)
	} else {
		ratio2 := float64(w_target) / float64(h_target)
		if ratio2 < ratio {
			w_resize = w_target
			h_resize = int(float64(w_target) / ratio)
		} else {
			h_resize = h_target
			w_resize = int(float64(h_target) * ratio)
		}
	}
	return w_resize, h_resize
}
