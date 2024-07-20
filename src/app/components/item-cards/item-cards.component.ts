import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { CdkAccordionModule } from '@angular/cdk/accordion';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { map, Observable, Subscription, switchMap, tap, timer } from 'rxjs';
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
  loading: boolean = false;
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
    // Due to the implementation, this timeout is necessary. Better implementation: Subscribe in a service.
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
      if (filters.dataType === 'flip' && item.margin < filters.margin.filter) {
        return false;
      } else if (
        this.filters.dataType === 'alch' &&
        item.alchprof! < filters.alchprof.filter
      ) {
        return false;
      }
      if (item.buylimit < this.filters.buyLimit.filter) return false;
      if (item.highpricevolume < this.filters.highVolume.filter) return false;
      if (item.lowpricevolume < this.filters.lowVolume.filter) return false;
      if (item.members && !this.filters.members) return false;
      return true;
    });

    // Sort data based on view
    if (this.filters.dataType === 'flip') {
      this.itemData.sort((a, b) => b.margin - a.margin);
    } else if (this.filters.dataType === 'alch') {
      this.itemData.sort((a, b) => b.alchprof! - a.alchprof!);
    }

    this.setMaxFilterValues();
  }

  setMaxFilterValues(): void {
    for (const item of this.itemData) {
      this.filters.alchprof.max = Math.max(
        this.filters.alchprof.max,
        item.alchprof!
      );
      this.filters.buyLimit.max = Math.max(
        this.filters.buyLimit.max,
        item.buylimit
      );
      this.filters.margin.max = Math.max(this.filters.margin.max, item.margin);
      this.filters.highVolume.max = Math.max(
        this.filters.highVolume.max,
        item.highpricevolume
      );
      this.filters.lowVolume.max = Math.max(
        this.filters.lowVolume.max,
        item.lowpricevolume
      );
    }
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
    this.loading = true;

    if (event.dataType !== this.filters.dataType) {
      this.filters.dataType = event.dataType;
    } else {
      this.filters = { ...event };
    }

    const applyFilters = () => {
      this.filterItems(this.filters);
    };

    requestAnimationFrame(applyFilters);
    setTimeout(() => (this.loading = false), 0);
  }
}
