import { Component, OnDestroy, OnInit } from '@angular/core';
import { FetchmarketdataService } from '../../services/fetchmarketdata.service';
import {
  firstValueFrom,
  map,
  Observable,
  of,
  Subscription,
  switchMap,
  tap,
  timer,
} from 'rxjs';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { CommonModule } from '@angular/common';
import { tick } from '@angular/core/testing';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  providers: [FetchmarketdataService],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent implements OnInit, OnDestroy {
  spreadData$: Observable<IItemListings[]> = new Observable<IItemListings[]>();
  showdata: IItemListings[] = [];
  subscription: Subscription = new Subscription();

  constructor(private _fetchMarketData: FetchmarketdataService) {}

  ngOnInit() {
    this.pollService();
  }

  ngOnDestroy(): void {
    if (this.subscription) this.subscription.unsubscribe();
  }

  pollService() {
    this.spreadData$ = this._fetchMarketData.getMarketData();
  }

  sub(): void {
    if (this.subscription) this.subscription.unsubscribe();

    this.subscription = this.spreadData$.subscribe({
      next: (data) => {
        this.showdata = data;
      },
      error: (error) => {
        console.error('Error subscribing to Observer: ', error);
      },
    });
  }
}
