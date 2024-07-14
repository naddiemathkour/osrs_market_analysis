export interface IFilters {
  dataType: string;
  margin: { max?: number; filter: number };
  alchprof: { max?: number; filter: number };
  buyLimit: { max?: number; filter: number };
  highVolume: { max?: number; filter: number };
  lowVolume: { max?: number; filter: number };
  members: boolean;
}
