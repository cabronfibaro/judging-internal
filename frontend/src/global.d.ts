import type { repository } from "$wails/go/models";

declare global {
  namespace Tournament {
    type GroupedMatches = repository.Match[][][];
  }
}

export {};
