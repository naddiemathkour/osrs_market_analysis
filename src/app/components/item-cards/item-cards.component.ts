import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { CdkAccordionModule } from '@angular/cdk/accordion';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { map, Observable, of, Subscription, switchMap, tap, timer } from 'rxjs';
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
  rawItemData: IItemListings[] = [];
  testBool: boolean = true;
  filters: IFilters = {} as IFilters;
  nat: IItemListings | undefined = {} as IItemListings;

  constructor(private _fetchMarketData: FetchmarketdataService) {}

  ngOnInit(): void {
    this.filters = {
      dataType: '',
      alchprof: { max: 0, filter: 1 },
      margin: { max: 0, filter: 5 },
      buyLimit: { max: 0, filter: 1 },
      highVolume: { max: 0, filter: 1 },
      lowVolume: { max: 0, filter: 1 },
      members: true,
    };
    this.fetchData();
    this.subscribeToData();
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
          switchMap(() =>
            this._fetchMarketData.getMarketData().pipe(
              tap(
                (data: IItemListings[]) =>
                  (this.nat = data.find((item) => item.name === 'Nature rune'))
              ),
              map((data: IItemListings[]) => {
                for (const item of data) {
                  item.alchprof =
                    item.highalch - item.avghighprice - this.nat!.avghighprice;
                  item.totalvolume = item.highpricevolume + item.lowpricevolume;
                  item.buylimitprof = item.alchprof * item.buylimit;
                }
                return data;
              })
            )
          )
        )),
      0
    );
  }

  subscribeToData() {
    if (this.marketData$)
      this._subscription = this.marketData$.subscribe((data) => {
        this.rawItemData = data;
        this.filterItems(this.filters);
      });
  }

  filterItems(filters: IFilters): void {
    // Filter through itemData based on filter options
    this.itemData = this.rawItemData.filter((item) => {
      if (
        filters.dataType === 'flip' &&
        item.margin < (filters.margin?.filter || 1)
      ) {
        return false;
      } else if (
        this.filters.dataType === 'alch' &&
        item.alchprof! - this.nat!.avghighprice < (filters.alchprof.filter || 1)
      ) {
        return false;
      }
      if (item.buylimit < (this.filters.buyLimit.filter || 1)) return false;
      if (item.highpricevolume < (this.filters.highVolume.filter || 1))
        return false;
      if (item.lowpricevolume < (this.filters.lowVolume.filter || 1))
        return false;
      if (item.members && !this.filters.members) return false;
      return true;
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

  updateFilters(event: IFilters) {
    if (event.dataType !== this.filters.dataType) {
      if (event.dataType === 'alch') event.margin.filter = 5;
      else event.alchprof.filter = 0;
    }
    this.filters = { ...event };
    this.filterItems(this.filters);
    console.log(this.filters);
  }
}
