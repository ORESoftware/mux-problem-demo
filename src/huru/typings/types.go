package typings

type Entities struct {
	basic struct {
		path string
		req  struct {
			headers struct {
				x_requested_by string
			}
			body struct {
				foo  string
				bar  int
				zoom bool
			}
		}
		res struct {
			headers struct {
			}
		}
	}
	tragic struct {
		path string
		req  struct {
			headers struct {
				x_requested_by string
			}
			body struct {
				foo string
			}
		}
		res struct {
			headers struct {
			}
		}
	}
}
