import { Component, OnInit } from '@angular/core';
import { MarketDataService } from '../../services/market-data.service';
import { HttpClientModule } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { ILatestData } from '../../interfaces/latest-data-interface';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [HttpClientModule, CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent implements OnInit {
  marketData$!: Observable<ILatestData[]>;
  latestData: ILatestData[] = [];
  itemIds: any[] = [];

  constructor(private marketDataService: MarketDataService) {}

  ngOnInit(): void {
    this.marketData$ = this.marketDataService.fetchMarketData();
    this.marketData$.subscribe((data) => {
      const t = data['data' as any];
      this.latestData = Object.values(t);
      this.itemIds = Object.keys(t);
      for (let i = 0; i < this.itemIds.length; i++) {
        this.latestData[i].id = this.itemIds[i];
      }
    })
  }
}