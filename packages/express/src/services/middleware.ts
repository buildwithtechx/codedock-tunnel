import type { RequestHandler } from 'express';
import type { ExpressTunnel } from '../interfaces/options';

export function tunnelStatus(tunnel: ExpressTunnel): RequestHandler {
  return (_request, response) => {
    response.json(tunnel.state());
  };
}

export function tunnelLifecycle(tunnel: ExpressTunnel): RequestHandler {
  return (_request, _response, next) => {
    void tunnel
      .start()
      .then(() => next())
      .catch(next);
  };
}
