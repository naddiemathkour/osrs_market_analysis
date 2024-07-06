import { Component } from '@angular/core';
import { FetchmarketdataService } from '../../services/fetchmarketdata.service';
import { map, Observable } from 'rxjs';
import { IItemListings } from '../../interfaces/itemlistings.interface';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [],
  providers: [FetchmarketdataService],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {
  spreadData: Observable<IItemListings[]> = [] as any;
  spreads: IItemListings[] = [];

  constructor(private _fetchMarketData: FetchmarketdataService) {}

  fetch() {
    this._fetchMarketData
      .getMarketData()
      .subscribe((data) => (this.spreads = data.items));
  }
}
