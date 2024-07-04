import { Component } from '@angular/core';
import { FetchmarketdataService } from '../../services/fetchmarketdata.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [],
  providers: [FetchmarketdataService],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {
  constructor(private _fetchMarketData: FetchmarketdataService) {}

  fetch() {
    this._fetchMarketData
      .getMarketData()
      .subscribe((data) => console.log(data));
  }
}
