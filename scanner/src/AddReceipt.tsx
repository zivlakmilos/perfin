import { useNavigate, useSearchParams } from '@solidjs/router';
import { createSignal, For, Show, type Component } from 'solid-js';
import { processReceipt, TReceipt } from './utils/api';
import { loadToken, token } from './utils/token';

const AddReceipt: Component = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [receipt, setReceipt] = createSignal<TReceipt | null>(null);

  loadToken();
  if (!token()) {
    navigate("/auth/login", { replace: true });
  }

  (async () => {
    if (searchParams.value) {
      try {
        const value = decodeURIComponent(searchParams.value as string);
        const res = await processReceipt(value);
        const rec: TReceipt = {
          taxId: res.receipt.invoiceRequest.taxId,
          businessName: res.receipt.invoiceRequest.locationName,
          date: res.receipt.invoiceResult.sdcTime as string,
          totalAmount: res.receipt.invoiceResult.totalAmount,
          paymentAccount: res.receipt.invoiceRequest.payments[0].paymentType.toString(), // TODO: Get account from mappings
          url: value,
          items: []
        };

        for (const item of res.items) {
          console.log(item);
          rec.items.push({
            name: item.title,
            price: item.price,
            quantity: item.quantity,
            amount: item.amount,
            account: '', // TODO: Get account from mappings
          });
        }

        console.log(rec);
        setReceipt(rec);
      } catch (err) {
        console.log(err);
        alert(err);
      }
    }
  })();

  return (
    <div class="px-2 pt-4 max-w-full w-full mx-auto">
      <Show when={receipt()}>
        {(rec) => (
          <form class="space-y-6">
            <div class="bg-base-200 w-full p-2 rounded-box">
              <div class="text-xl font-medium mb-2">
                Receipt Info
              </div>
              <div class="space-y-4">
                <div class="form-control">
                  <label class="label block mb-1">
                    <span class="label-text">Tax ID</span>
                  </label>
                  <input class="input input-bordered w-full" type="text" value={rec().taxId} readonly />
                </div>
                <div class="form-control">
                  <label class="label block mb-1">
                    <span class="label-text">Business Name</span>
                  </label>
                  <input class="input input-bordered w-full" type="text" value={rec().businessName} readonly />
                </div>
                <div class="form-control">
                  <label class="label block mb-1">
                    <span class="label-text">Date</span>
                  </label>
                  <input class="input input-bordered w-full" type="text" value={rec().date} readonly />
                </div>
                <div class="form-control">
                  <label class="label block mb-1">
                    <span class="label-text">Total Amount</span>
                  </label>
                  <input class="input input-bordered w-full" type="number" value={rec().totalAmount} readonly />
                </div>
                <div class="form-control">
                  <label class="label block mb-1">
                    <span class="label-text">Payment Account</span>
                  </label>
                  <input class="input input-bordered w-full" type="text" value={rec().paymentAccount} readonly />
                </div>
                <div class="form-control">
                  <label class="label block mb-1">
                    <span class="label-text">Receipt URL</span>
                  </label>
                  <input class="input input-bordered w-full" type="text" value={rec().url} readonly />
                </div>
              </div>
            </div>
            <div class="bg-base-200 w-full p-2 rounded-box">
              <div class="text-xl font-medium mb-2">
                Items
              </div>
              <div>
                <div class="overflow-x-auto">
                  <table class="table table-zebra">
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Account</th>
                      </tr>
                    </thead>
                    <tbody>
                      <For each={rec().items}>
                        {(item, idx) => (
                          <tr>
                            <td>
                              <input
                                class="input input-bordered input-sm"
                                type="text"
                                value={item.name}
                                onInput={e => {
                                  const items = [...rec().items];
                                  items[idx()].name = e.currentTarget.value;
                                  setReceipt({ ...rec(), items });
                                }}
                              />
                            </td>
                            <td>
                              <input
                                class="input input-bordered input-sm"
                                type="text"
                                value={item.account}
                                onInput={e => {
                                  const items = [...rec().items];
                                  items[idx()].account = e.currentTarget.value;
                                  setReceipt({ ...rec(), items });
                                }}
                              />
                            </td>
                          </tr>
                        )}
                      </For>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            <div class="form-control mt-4">
              <button class="btn btn-primary" type="submit">Save Receipt</button>
            </div>
          </form>
        )}
      </Show>
    </div>
  );
}

export default AddReceipt;
