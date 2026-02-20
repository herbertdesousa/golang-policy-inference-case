import http from "k6/http";
import { check, sleep } from "k6";

import { mapPayloadExpected } from "./payloads.ts";
import { deepEqual, pickRandom } from "./utils.ts";

export const options = {
  thresholds: {
    http_req_duration: ["p(99) < 500"],
  },
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 50,
      timeUnit: "1s",
      duration: "60s",
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
};

export default function () {
  const url = "http://localhost:8080/infer";
  const { payload } = pickRandom(mapPayloadExpected);
  const payloadString = JSON.stringify(payload);
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const res = http.post(url, payloadString, params);

  check(res, {
    "status was 200": (r) => r.status === 200,
  });

  sleep(1);
}
