import { useNavigate, useSearchParams } from '@solidjs/router';
import type { Component } from 'solid-js';

const AddReceipt: Component = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  let value = "";
  if (searchParams.value) {
    value = decodeURIComponent(searchParams.value as string);
  }

  return (
    <div>
      {value}
    </div>
  );
};

export default AddReceipt;
