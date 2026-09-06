import React from 'react';

interface Transaction {
  id: string;
  amount: number;
}

export const TransactionList = ({ items }: { items: Transaction[] }) => {
  return (
    <ul>
      {items.map((item) => (
        <li key={item.id}>
          <span>{item.id}</span>
          <span>{item.amount}</span>
        </li>
      ))}
    </ul>
  );
};
