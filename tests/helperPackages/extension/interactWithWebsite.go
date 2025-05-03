package extension

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

func (ce *ChromeExtension) MakeThePaymentAndGetOnPaidTier(ctx context.Context, selectOneTimeButton bool) error {
	ctx, _ = context.WithTimeout(ctx, 700*time.Second)
	// not calling the cancel function as if we do then the future tests will not work as they get cancelled
	// defer cancel()
	//

	err := ce.NavigateToAWebPage(ctx, "http://localhost:5173/")
	if err != nil {
		println("there is a error in navigatign to the website", err.Error())
		return err
	}
	println("we have navigated to the website and now we will sleeep for 7 sec and then wait for the button to be visible and then click the button")
	time.Sleep(3 * time.Second)
	tryOnceButton := `bg-yellow-500 text-black flex w-full transform items-center justify-center rounded-md border-3 border-black px-8 py-3 font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-[4px] active:translate-y-[4px] active:shadow-none s-N59DAw_bCMv0`
	goPremiumButton := `bg-purple-600 text-white flex w-full transform items-center justify-center rounded-md border-3 border-black px-8 py-3 font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-[4px] active:translate-y-[4px] active:shadow-none s-N59DAw_bCMv0`
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(tryOnceButton, chromedp.BySearch),
	)
	println("we have waited for the button to be visible and now we will click on it, hopefully it is visible")
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	println("1")
	println("we are about to click on the button")
	selectedButton := ""
	if selectOneTimeButton {
		selectedButton = tryOnceButton
	} else {
		selectedButton = goPremiumButton
	}
	err = chromedp.Run(ctx,
		chromedp.Click(selectedButton, chromedp.BySearch, chromedp.NodeVisible),
	)
	if err != nil {
		return err
	}
	println("about to click via nodes")
	println("we are going to click on the bank")
	time.Sleep(time.Second * 4)
	println("hopefully working")
	time.Sleep(time.Second * 2)
	// document.querySelector("#nav-sidebar > div:nth-child(1) > label:nth-child(2) > div > div")
	var iframes []*cdp.Node
	err = chromedp.Run(ctx, chromedp.Nodes("iframe", &iframes, chromedp.ByQuery))
	if err != nil {
		println("error in getting the iframes ->", err.Error(), " and the iframes are ->", len(iframes))
		return err
	}

	var iframeT *target.Info
	targets, err := chromedp.Targets(ctx)

	println("the targets are ->", len(targets))
	if err != nil {
		println("error in getting the targets and it is ->", err.Error())
		return err
	}
	for _, t := range targets {
		println("targets tyoe is ->", t.Type, " and the url is ->", t.URL)
		if t.Type == "iframe" && strings.Contains(t.URL, "razorpay.com") {
			iframeT = t
			break
		}
	}
	if iframeT == nil {
		return fmt.Errorf("iframe target not found")
	}

	iframeCtx, _ := chromedp.NewContext(
		ctx,
		chromedp.WithTargetID(iframeT.TargetID),
	)
	// not calling the cancel function as if we do then the future tests will not work as they get cancelled
	// defer iframeCancel()

	// var fullHTML string
	// err = chromedp.Run(iframeCtx,
	// 	chromedp.Sleep(time.Second*2),
	// 	chromedp.Evaluate(`document.documentElement.outerHTML`, &fullHTML),
	// ) //
	// if err != nil {
	// 	println("error:", err.Error())
	// }
	// println("the full html  form the  iframe is ->\n ->", fullHTML)

	println("about to click on the button form the iframe ctx")

	err = chromedp.Run(iframeCtx,
		// for example, click the Netbanking button

		chromedp.Click(`//*[@id="nav-sidebar"]/div[1]/label[2]/div/div`, chromedp.BySearch),
	)
	if err != nil {
		return err
	}

	// this is here to get the new window that the razorpay opens
	println("clicking on the bank of borada in 3 sec")
	time.Sleep(time.Second * 3)
	err = chromedp.Run(iframeCtx,
		chromedp.Click(`//*[@id="main-stack-container"]/div/div/div/div/div[2]/div/div/form[1]/div/label[1]/div/div`, chromedp.BySearch),
	)
	println("clicked the bank of borada button")
	if err != nil {
		println("there is an error in clicking the bank of borada button and it is ->", err.Error())
		return err
	}

	time.Sleep(time.Second * 5)
	println("finding the new targets")

	var successPageTarget *target.Info
	targets, err = chromedp.Targets(ctx)

	println("the targets are ->", len(targets))
	if err != nil {
		println("error in getting the targets and it is ->", err.Error())
		return err
	}

	for _, t := range targets {
		println("target's type is ->", t.Type, " and the url is ->", t.URL, " and the title is ", t.Title, "==\n")
	}
	for _, t := range targets {
		// println("target's type is ->", t.Type, " and the url is ->", t.URL, " and the title is ", t.Title)
		if t.Type == "page" && strings.Contains(t.URL, "razorpay.com") && strings.Contains(t.Title, "Razorpay Bank") {
			println("-- selected target's type is ->", t.Type, " and the url is ->", t.URL, " and the title is ", t.Title)
			successPageTarget = t
			break
		}
	}
	if successPageTarget == nil {
		return fmt.Errorf("successPage target not found")
	}

	successPageCtx, _ := chromedp.NewContext(
		iframeCtx,
		chromedp.WithTargetID(successPageTarget.TargetID),
	)

	// not calling the cancel function as if we do then the future tests will not work as they get cancelled
	// defer successPageCancelFunc()

	// var fullHTML string
	// err = chromedp.Run(successPageCtx,
	// 	chromedp.Sleep(time.Second*2),
	// 	chromedp.Evaluate(`document.documentElement.outerHTML`, &fullHTML),
	// ) //
	// if err != nil {
	// 	println("error:", err.Error())
	// }
	// println("the full html  form the  success page  is ->\n ->", fullHTML)
	time.Sleep(time.Second * 3)

	println("clicking on the success button")

	err = chromedp.Run(successPageCtx,
		// chromedp.WaitVisible(`/html/body/form/button[1]`, chromedp.BySearch),
		chromedp.Click(`/html/body/form/button[1]`, chromedp.BySearch),
	)
	println("clicked the success button")
	if err != nil {
		println("there is an error in clicking the succcess button and it is ->", err.Error())
		return err
	}
	println("sleeping for 24 sec")
	time.Sleep(time.Second * 24)
	println("we are able to click the success button and now we will now return ")
	return nil
	//	println("we got the iframes and it is ->", len(iframes))
	//
	//	fmt.Printf("the iframe struct is -> %+v \n\n", iframes[0])
	//	println("clikcing on the netbanking button")
	//	println("possible iframe id is ->", iframes[0].FrameID)
	//	// iframeID := iframes[0].FrameID
	//	b := fmt.Sprintf(`
	//
	// let a = document.querySelector("body > div.razorpay-container > iframe:nth-child(3)").contentWindow.document.querySelector('#nav-sidebar > div:nth-child(1) > label:nth-child(2) > div > div')
	//
	//	  console.log('-------pp----------'+a)
	//	  a.click();
	//	    `)
	//		err = chromedp.Run(ctx, chromedp.Evaluate(b, nil))
	//		if err != nil {
	//			println("there is a error in clicking the button with eval the selected frames and it is ->", err.Error())
	//			time.Sleep(time.Second * 15)
	//			return err
	//		}
	//		err = chromedp.Run(ctx,
	//			chromedp.Click(
	//				`//*[@id="nav-sidebar"]/div[1]/label[2]/div/div`,
	//				chromedp.BySearch,             // use CSS selector
	//				chromedp.FromNode(iframes[0]), // scope into that iframe node&#8203;:contentReference[oaicite:0]{index=0}
	//			),
	//		)
	//		// err = clickViaIframeContext(ctx)
	//		if err != nil {
	//			println("there is a error in clicking the button wiht the selected frames and it is ->", err.Error())
	//			return err
	//		}
	//		println("out iof it ")
	//		// err = chromedp.Run(ctx, chromedp.Click(`*[@id="nav-sidebar"]/div[1]/label[2]/div/div`, chromedp.BySearch, chromedp.FromNode(iframes[0])))
	//		//#nav-sidebar > div:nth-child(1) > label:nth-child(2) > div > div:nth-child(1)
	//		err = chromedp.Run(ctx, chromedp.Click(`#nav-sidebar > div:nth-child(1) > label:nth-child(2) > div > div:nth-child(1)`, chromedp.ByQuery, chromedp.FromNode(iframes[0])))
	//		if err != nil {
	//			println("error in clicking the button using the nodes  ->", err.Error())
	//			return err
	//		}
	//		println("will continue with the clicking after 15 sec")
	//		time.Sleep(time.Second * 15)
	//		println("continuing ...")
	//		err = chromedp.Run(ctx,
	//			chromedp.EvaluateAsDevTools(`(function(){
	//
	// const btn = document.querySelector("#nav-sidebar > div:nth-child(1) > label:nth-child(2) > div > div")
	//
	//	      console.log("\n\n\n\n +++the button is ->", btn)
	//	      if (btn === null) {
	//	        return JSON.stringify({ success: false, error: "button was  not found:-> "+btn });
	//	      }
	//	      btn.click();
	//	      return JSON.stringify({ success: true, text: btn });
	//	  })()`, &result),
	//			chromedp.Sleep(time.Second*100),
	//		)
	//		if err != nil {
	//
	//			println("there is an error in clicking the netbanking button and it is ->", err.Error())
	//			time.Sleep(time.Second * 30)
	//			return err
	//		}
	//		err = chromedp.Run(ctx,
	//			chromedp.Evaluate(`
	//	      let a = document.querySelector("#main-stack-container > div > div > div > div > div.flex.h-full.flex-1.flex-col.overflow-auto.bg-surface > div > div > form:nth-child(4) > div > label:nth-child(1) > div > div > div > div")
	//	      a.click()
	//	      `, nil),
	//			chromedp.Sleep(time.Second*2),
	//		)
	//		println("the result form clicking the button is ->", result)
	//		if err != nil {
	//			return err
	//		}
	//		println("we have clicked the button and now we are waiting for the payment to be done, we will sleep for next 100 seconds")
	//		time.Sleep(time.Second * 100)
}
