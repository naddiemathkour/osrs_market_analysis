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
    this.fetchData();
  }

  ngOnDestroy(): void {
    this._subscription.unsubscribe();
  }

  fetchData(): void {
    // Poll http service call to refresh data every minute
    this.marketData$ = timer(0, 1 * 60 * 1000).pipe(
      switchMap(() => this._fetchMarketData.getMarketData())
    );
  }

  testFunc(): void {
    console.log('Test func');
  }

  filterItems(filters: IFilters): void {
    this._subscription = this.marketData$.subscribe((data) => {
      this.nat = data.filter((item) => item.name === 'Nature rune')[0];
      if (filters.dataType === 'flip') {
        this.itemData = data.filter(
          (item) => item.margin > (filters.margin || 1)
        );
      } else if (filters.dataType === 'alch') {
        this.itemData = data.filter(
          (item) =>
            item.highalch - item.avghighprice - this.nat.avghighprice >
            (filters.alchprof || 1)
        );
      }
    });

    this.itemData = this.itemData.filter((item) => {
      if (!filters.members && item.members) return false;
      if (item.buylimit < (filters.buyLimit || 0)) return false;
      if (item.highpricevolume < (filters.highVolume || 1)) return false;
      if (item.lowpricevolume < (filters.lowVolume || 1)) return false;
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
}
