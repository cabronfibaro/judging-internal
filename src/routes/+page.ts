import type { Tag } from '../types';

type Category = {
  label: string;
  tags: Tag[];
};

type CriteriaConfig = {
  name: string;
  categories: Record<string, Category>;
};

export async function load({ fetch }) {
  const response = await fetch("/conf/default.json");

  if (!response.ok) {
    throw new Error("Failed to load data.json");
  }

  const criteriaConfig = (await response.json()) as CriteriaConfig;

  return {
    criteriaConfig,
  };
}
