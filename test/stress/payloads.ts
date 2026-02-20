type PayloadExpected = {
  payload: {
    policy_dot: Window | null;
    input: Record<string, any>;
  };
};

const dotComplexPolicy = open("../digraphs/complex_policy.dot");
const dotAutoInsurance = open("../digraphs/auto_insurance.dot");
const dotB2bLoan = open("../digraphs/b2b_loan_policy.dot");
const dotEcommerceFraud = open("../digraphs/ecommerce_fraud.dot");

export const mapPayloadExpected: PayloadExpected[] = [
  {
    payload: {
      policy_dot: dotComplexPolicy,
      input: {
        age: 17,
        score: 0,
        income: 0,
        country: "BR",
        tenure_months: 0,
        risk_flag: false,
        has_default: false,
        product_type: "credit",
      },
    },
  },
  {
    payload: {
      policy_dot: dotComplexPolicy,
      input: {
        age: 30,
        score: 800,
        income: 20000,
        country: "BR",
        tenure_months: 48,
        risk_flag: false,
        has_default: false,
        product_type: "credit",
      },
    },
  },
  {
    payload: {
      policy_dot: dotAutoInsurance,
      input: {
        driver_age: 15,
        accidents_5y: 0,
        tickets_3y: 0,
        vehicle_value: 0,
        credit_score: 0,
      },
    },
  },
  {
    payload: {
      policy_dot: dotAutoInsurance,
      input: {
        driver_age: 30,
        accidents_5y: 0,
        tickets_3y: 0,
        vehicle_value: 30000,
        credit_score: 800,
      },
    },
  },
  {
    payload: {
      policy_dot: dotAutoInsurance,
      input: {
        driver_age: 18,
        accidents_5y: 2,
        tickets_3y: 1,
        vehicle_value: 15000,
        credit_score: 650,
      },
    },
  },
  {
    payload: {
      policy_dot: dotAutoInsurance,
      input: {
        driver_age: 22,
        accidents_5y: 1,
        tickets_3y: 1,
        vehicle_value: 90000,
        credit_score: 650,
      },
    },
  },
  {
    payload: {
      policy_dot: dotAutoInsurance,
      input: {
        driver_age: 70,
        accidents_5y: 2,
        tickets_3y: 0,
        vehicle_value: 20000,
        credit_score: 750,
      },
    },
  },
  {
    payload: {
      policy_dot: dotB2bLoan,
      input: {
        has_recent_bankruptcy: true,
        company_age_years: 0,
        annual_revenue: 0,
        current_debt: 0,
        owner_credit_score: 0,
        industry_risk: "high",
      },
    },
  },
  {
    payload: {
      policy_dot: dotB2bLoan,
      input: {
        has_recent_bankruptcy: false,
        company_age_years: 5,
        annual_revenue: 25000000,
        current_debt: 1000000,
        owner_credit_score: 800,
        industry_risk: "low",
      },
    },
  },
  {
    payload: {
      policy_dot: dotB2bLoan,
      input: {
        has_recent_bankruptcy: false,
        company_age_years: 2,
        annual_revenue: 1000000,
        current_debt: 600000,
        owner_credit_score: 750,
        industry_risk: "low",
      },
    },
  },
  {
    payload: {
      policy_dot: dotB2bLoan,
      input: {
        has_recent_bankruptcy: false,
        company_age_years: 10,
        annual_revenue: 1500000,
        current_debt: 200000,
        owner_credit_score: 650,
        industry_risk: "high",
      },
    },
  },
  {
    payload: {
      policy_dot: dotEcommerceFraud,
      input: {
        account_age_days: 0,
        total_past_purchases: 0,
        orders_24h: 3,
        cart_value: 50,
        ip_country: "US",
        shipping_country: "US",
        is_vpn: false,
        device_trust_score: 100,
      },
    },
  },
  {
    payload: {
      policy_dot: dotEcommerceFraud,
      input: {
        account_age_days: 100,
        total_past_purchases: 15,
        orders_24h: 1,
        cart_value: 500,
        ip_country: "US",
        shipping_country: "US",
        is_vpn: false,
        device_trust_score: 95,
      },
    },
  },
  {
    payload: {
      policy_dot: dotEcommerceFraud,
      input: {
        account_age_days: 10,
        total_past_purchases: 0,
        orders_24h: 1,
        cart_value: 100,
        ip_country: "RU",
        shipping_country: "US",
        is_vpn: true,
        device_trust_score: 50,
      },
    },
  },
  {
    payload: {
      policy_dot: dotEcommerceFraud,
      input: {
        account_age_days: 40,
        total_past_purchases: 5,
        orders_24h: 1,
        cart_value: 200,
        ip_country: "CA",
        shipping_country: "CA",
        is_vpn: false,
        device_trust_score: 65,
      },
    },
  },
];
