import { json, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = ({ request }) => {
  console.log(...request.headers);

  return json({
    message: "OK!",
  });
};
