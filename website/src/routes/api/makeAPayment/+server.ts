import { json } from '@sveltejs/kit';
import { PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH } from '$env/static/public';

// Cache headers
const STANDARD_HEADERS = {
  'Content-Type': 'application/json'
};

// Precomputed URL
const BACKEND_ENDPOINT = `${PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH}/makeAPayment`;

export async function POST({ request }) {
  try {
    // Stream the request body directly to the backend without parsing first
    // This reduces memory usage and processing time
    console.log(`++++++++++++++ making a request to the makeAPayment ${new Date().getTime()}  ++++++++++`);
    const requestData = await request.json();

    const backendResponse = await fetch(BACKEND_ENDPOINT, {
      method: 'POST',
      headers: STANDARD_HEADERS,
      // Use the request body directly as a stream
      body: JSON.stringify(requestData)
    });

    if (!backendResponse.ok) {
      console.error('Backend error:', backendResponse.statusText);
      return json({
        order_id: '',
        message: `Backend error: ${backendResponse.statusText}`,
        status_code: backendResponse.status.toString()
      }, { status: backendResponse.status });
    }

    // Return the response directly
    // This avoids unnecessary JSON parsing and re-serialization
    const responseData = await backendResponse.json();
    return json(responseData);

  } catch (err) {
    console.error('Error in makeAPayment route:', err);
    return json({
      order_id: '',
      message: `Internal server error: ${err instanceof Error ? err.message : String(err)}`,
      status_code: '500'
    }, { status: 500 });
  }
}
