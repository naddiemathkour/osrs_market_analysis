export interface IItemListings {
  id: number;
  name: string;
  icon: string;
  examine: string;
  members: true;
  buylimit: number;
  highalch: number;
  lowalch: number;
  avghighprice: number;
  highpricevolume: number;
  avglowprice: number;
  lowpricevolume: number;
  spread: number;
  margin: number;
  timestamp: Date;
}
