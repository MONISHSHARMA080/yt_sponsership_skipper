import { json, type RequestHandler } from '@sveltejs/kit';
import { PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH } from "$env/static/public";

export const POST: RequestHandler = async ({ request }) => {
  try {
    // Get the payment data from the request body
    const paymentData = await request.json();

    // Forward the request to the Go backend
    const response = await fetch(`${PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH}/validatePayment`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(paymentData)
    });

    if (!response.ok) {
      throw new Error(`Go backend returned status: ${response.status}`);
    }
    console.log("the response data is ", paymentData);

    // Get the response data from the Go backend
    const data = await response.json();

    // Return the same data structure from SvelteKit backend
    return json(data);
  } catch (error) {
    console.error('Error processing payment validation:', error);
    return json({
      message: error instanceof Error ? error.message : 'Unknown error',
      status_code: 500,
      success: false
    }, { status: 500 });
  }
};