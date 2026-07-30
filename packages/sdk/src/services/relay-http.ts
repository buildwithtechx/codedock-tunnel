import type { HTTPRequest, HTTPResponse } from '../protocol';

type ResponseSender = (response: HTTPResponse) => void;

export async function forwardHttpRequest(
  request: HTTPRequest,
  requestId: string | undefined,
  localPort: number | undefined,
  send: ResponseSender,
): Promise<void> {
  if (!requestId) return;
  if (!localPort || typeof fetch !== 'function') {
    sendSafely(send, {
      status_code: 502,
      headers: {},
      error: 'local HTTP forwarding is unavailable',
    });
    return;
  }
  try {
    const headers = new Headers();
    for (const [name, values] of Object.entries(request.headers ?? {})) {
      headers.set(name, values.join(', '));
    }
    const response = await fetch(
      `http://127.0.0.1:${localPort}${request.path}`,
      {
        method: request.method,
        headers,
        body: request.body
          ? (decodeBase64(request.body).buffer as ArrayBuffer)
          : undefined,
      },
    );
    const responseHeaders: Record<string, string[]> = {};
    response.headers.forEach((value, name) => {
      responseHeaders[name] = [value];
    });
    sendSafely(send, {
      status_code: response.status,
      headers: responseHeaders,
      body: encodeBase64(new Uint8Array(await response.arrayBuffer())),
    });
  } catch (error) {
    sendSafely(send, {
      status_code: 502,
      headers: {},
      error: error instanceof Error ? error.message : String(error),
    });
  }
}

function sendSafely(send: ResponseSender, response: HTTPResponse): void {
  try {
    send(response);
  } catch {
    return;
  }
}

function decodeBase64(value: string): Uint8Array {
  const binary = atob(value);
  return Uint8Array.from(binary, (character) => character.charCodeAt(0));
}

function encodeBase64(value: Uint8Array): string {
  let binary = '';
  for (const byte of value) binary += String.fromCharCode(byte);
  return btoa(binary);
}
