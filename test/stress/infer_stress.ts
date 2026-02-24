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
    ramping_request_rate: {
      executor: "ramping-arrival-rate",
      startRate: 50,
      timeUnit: "1s",
      preAllocatedVUs: 50,
      maxVUs: 500,
      stages: [
        { target: 50, duration: "30s" },
        { target: 100, duration: "10s" },
        { target: 200, duration: "30s" },
      ],
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
