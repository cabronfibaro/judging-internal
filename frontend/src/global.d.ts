import type { repository } from "$wails/go/models";

declare global {
  namespace Tournament {
    type GroupedMatches = repository.Match[][][];
  }

  namespace Judging {
    type Option = {
      label: string;
      value: number;
    };

    type Parameter = {
      name: string;
      type: string;
    };

    type Tag = {
      id: string;
      label: string;
      type: "bonus" | "deduction" | "neutral";
      value_type: "dynamic" | "enum" | "switch";
      options?: Option[];
      default_value?: number;
      parameter?: Parameter;
      max?: number;
      multiplier?: number;
    };

    type CriteriaCategory = {
      label: string;
      tags: Tag[];
    };

    type CriteriaConfig = {
      name: string;
      categories: Record<string, CriteriaCatrgory>;
    };
  }
}

export {};
