import { token } from "./token";

const BASE_URL = "http://localhost:9999"
//const BASE_URL = "http://192.168.0.101:9999"
//const BASE_URL = "http://finance.zivlak.rs/api"

export type TLogin = {
  token: string,
}

export type TFiscalReceiptInfo = {
  invoiceRequest: {
    posTime?: string;
    taxId: string;
    businessName: string;
    locationName: string;
    address: string;
    city: string;
    administrativeUnit: string;
    buyer?: string;
    buyerCostCenter?: string;
    cashier?: string;
    requestedBy: string;
    referentDocumentNumber?: string;
    invoiceType: number;
    transactionType: number;
    payments: Array<{
      paymentType: number;
      paymentTypeDescript: string;
      amount: number;
    }>;
  };
  invoiceResult: {
    totalAmount: number;
    transactionTypeCounter: number;
    totalCounter: number;
    invoiceCounterExtension: string;
    invoiceNumber: string;
    signedBy: string;
    sdcTime: string;
  };
  journal: string;
  isValid: boolean;
};

export type TFiscalReceiptItem = {
  title: string;
  price: number;
  quantity: number;
  amount: number;
};

export type TFiscalReceiptProcessResult = {
  items: TFiscalReceiptItem[];
  receipt: TFiscalReceiptInfo;
};

export type TReceipt = {
  taxId: string;
  businessName: string;
  date: string;
  totalAmount: number;
  paymentAccount: string;
  url: string;
  items: Array<{
    name: string;
    price: number;
    quantity: number;
    amount: number;
    account: string;
  }>;
};

const postRequest = async (url: string, body: any, token: string | null): Promise<any> => {
  const headers: any = {
    "Content-Type": "application/json",
  }
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    method: "POST",
    body: JSON.stringify(body),
    headers: headers,
  });
  if (response.status < 200 || response.status > 299) {
    try {
      const data = await response.json();
      throw data.error;
    } catch (err) {
      throw err;
    }
  }
  return response.json();
}

export const login = async (username: string, password: string): Promise<string> => {
  const data = await postRequest(`${BASE_URL}/auth/login`, {
    username: username,
    password: password,
  }, null) as TLogin;
  return data.token;
}

export const processReceipt = async (url: string): Promise<TFiscalReceiptProcessResult> => {
  const data = await postRequest(`${BASE_URL}/fiscal_receipts/process`, {
    receiptUrl: url,
  }, token()) as TFiscalReceiptProcessResult;
  return data;
}
