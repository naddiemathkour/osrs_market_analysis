import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { CdkAccordionModule } from '@angular/cdk/accordion';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { Observable, Subscription, switchMap, timer } from 'rxjs';
import { FetchmarketdataService } from '../../services/fetchmarketdata.service';
import { CommonModule } from '@angular/common';
import { IFilters } from '../../interfaces/filters.interface';
import { ListingFiltersComponent } from '../listing-filters/listing-filters.component';

@Component({
  selector: 'app-item-cards',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    CdkAccordionModule,
    ListingFiltersComponent,
  ],
  providers: [FetchmarketdataService],
  templateUrl: './item-cards.component.html',
  styleUrl: './item-cards.component.scss',
})
export class ItemCardsComponent implements OnInit, OnDestroy {
  private _subscription: Subscription = new Subscription();
  marketData$!: Observable<IItemListings[]>;
  accordionState: Map<number, boolean> = new Map();
  itemData: IItemListings[] = [];
  testBool: boolean = true;
  filters: IFilters = {} as IFilters;
  nat = {} as IItemListings;

  constructor(private _fetchMarketData: FetchmarketdataService) {}

  ngOnInit(): void {
    this.filters = {
      dataType: '',
      alchprof: { max: 0, filter: 1 },
      margin: { max: 0, filter: 1 },
      buyLimit: { max: 0, filter: 1 },
      highVolume: { max: 0, filter: 1 },
      lowVolume: { max: 0, filter: 1 },
      members: true,
    };
    this.fetchData();
  }

  ngOnDestroy(): void {
    this._subscription.unsubscribe();
  }

  fetchData(): void {
    // Poll http service call to refresh data every minute
    // I don't know why, but this needs to timeout for 0 milliseconds to work
    setTimeout(
      () =>
        (this.marketData$ = timer(0, 1 * 60 * 1000).pipe(
          switchMap(() => this._fetchMarketData.getMarketData())
        )),
      0
    );
  }

  subscribeToData() {
    if (this.marketData$)
      this._subscription = this.marketData$.subscribe((data) => {
        this.itemData = data;
        this.filterItems(this.filters);
      });
  }

  filterItems(filters: IFilters): void {
    console.log('Filters: ', filters);
    // Set nature rune value for alch pricing
    this.nat = this.itemData.filter((item) => item.name === 'Nature rune')[0];

    // Filter through itemData based on filter options
    this.itemData = this.itemData.filter((item) => {
      if (filters.dataType === 'flip') {
        return item.margin > (filters.margin?.filter || 1);
      } else {
        return (
          item.highalch - item.avghighprice - this.nat.avghighprice >
          (filters.alchprof?.filter || 1)
        );
      }
    });
  }

  setAccordionState(id: number): void {
    this.accordionState.has(id)
      ? this.accordionState.set(id, !this.accordionState.get(id))
      : this.accordionState.set(id, true);
  }

  getAccordionState(id: number): boolean {
    return this.accordionState.get(id) || false;
  }
}
