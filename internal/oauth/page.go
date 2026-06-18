package oauth

import (
	_ "embed"
	"fmt"
)

// edgedeltaLogoSVG is the Edge Delta wordmark shown on the browser result page.
//
//go:embed assets/edgedelta-logo.svg
var edgedeltaLogoSVG string

const (
	checkIcon = `<svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>`
	crossIcon = `<svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18M6 6l12 12"/></svg>`
)

// resultPage renders the branded HTML the browser lands on after the OAuth
// redirect. ok selects the success (green) vs failure (red) treatment.
func resultPage(ok bool, heading, message string) string {
	accent, icon := "#16a34a", checkIcon
	if !ok {
		accent, icon = "#dc2626", crossIcon
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="en"><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>edx — %[2]s</title>
<style>
  :root{--ink:#0b1f33;--muted:#5b6b7b;--line:#e7ebf0;--accent:#0a78e6;--c:%[1]s}
  *{box-sizing:border-box}
  html,body{height:100%%;margin:0}
  body{background:radial-gradient(1100px 520px at 50%% -8%%,#eaf2fb 0%%,#f6f9fc 46%%,#f6f9fc 100%%);
    color:var(--ink);font:15px/1.55 -apple-system,BlinkMacSystemFont,"Segoe UI",Inter,Roboto,Helvetica,Arial,sans-serif;
    -webkit-font-smoothing:antialiased;display:flex;align-items:center;justify-content:center;padding:24px}
  .card{width:min(420px,92vw);background:#fff;border:1px solid var(--line);border-radius:16px;padding:40px 36px;
    text-align:center;box-shadow:0 1px 2px rgba(11,31,51,.04),0 14px 36px rgba(11,31,51,.09)}
  .logo svg{height:26px;width:auto;display:block;margin:0 auto 28px}
  .badge{width:58px;height:58px;border-radius:50%%;display:inline-flex;align-items:center;justify-content:center;
    background:color-mix(in srgb,var(--c) 13%%,#fff);color:var(--c);margin-bottom:20px}
  h1{font-size:19px;font-weight:650;letter-spacing:-.01em;margin:0 0 8px}
  p{margin:0 auto;max-width:30ch;color:var(--muted);font-size:14px}
  .foot{margin-top:28px;padding-top:18px;border-top:1px solid var(--line);color:#9aa8b6;font-size:12px;letter-spacing:.01em}
  .foot b{color:var(--muted);font-weight:600}
</style></head>
<body>
  <main class="card">
    <div class="logo">%[3]s</div>
    <div class="badge">%[4]s</div>
    <h1>%[2]s</h1>
    <p>%[5]s</p>
    <div class="foot"><b>edx</b> &middot; Edge Delta CLI</div>
  </main>
</body></html>`, accent, heading, edgedeltaLogoSVG, icon, message)
}
