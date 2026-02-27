import http from "k6/http";
import { check } from "k6";
import { mapPayload } from "./payloads.ts";
import { pickRandom } from "./utils.ts";

export const options = {
  thresholds: {
    http_req_duration: ["p(90) < 30"],
    http_req_failed: ["rate<0.01"],
  },
  scenarios: {
    fixed_200_RPS: {
      executor: "constant-arrival-rate",
      rate: 200,
      duration: "30s",
      preAllocatedVUs: 100,
      maxVUs: 400,
    },
  },
};

export default function () {
  const url = "http://localhost:8080/infer";
  const { payload } = pickRandom(mapPayload);

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const res = http.post(url, JSON.stringify(payload), params);

  check(res, {
    "status was 200": (r) => r.status === 200,
  });
}
