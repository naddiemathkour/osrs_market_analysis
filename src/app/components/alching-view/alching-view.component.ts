import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subscription } from 'rxjs';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { CdkAccordionModule } from '@angular/cdk/accordion';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-alching-view',
  standalone: true,
  imports: [CommonModule, CdkAccordionModule],
  templateUrl: './alching-view.component.html',
  styleUrl: './alching-view.component.scss',
})
export class AlchingViewComponent implements OnInit, OnDestroy {
  @Input() listingData$!: Observable<IItemListings[]>;
  @Input() showMembers!: boolean;
  itemData: IItemListings[] = [];
  expandedItem: number = -1;
  accordionState: Map<number, boolean> = new Map();
  natureRune: IItemListings = {} as IItemListings;
  private _subscription: Subscription = new Subscription();

  ngOnInit(): void {
    this._subscription = this.listingData$.subscribe((data) => {
      this.natureRune = data.filter((item) => item.name === 'Nature rune')[0];
      this.itemData = data
        .filter(
          (a) =>
            a.highalch - a.avghighprice - this.natureRune.avghighprice > 0 &&
            a.lowpricevolume + a.highpricevolume >= 25
        )
        .sort(
          (a, b) =>
            b.highalch -
            b.avghighprice -
            this.natureRune.avghighprice -
            (a.highalch - a.avghighprice - this.natureRune.avghighprice)
        );
    });
  }

  ngOnDestroy(): void {
    this._subscription.unsubscribe();
  }

  setAccordionState(id: number): void {
    if (this.accordionState.has(id))
      this.accordionState.set(id, !this.accordionState.get(id));
    else this.accordionState.set(id, true);
  }

  getAccordionState(id: number): boolean {
    if (this.accordionState.has(id)) return this.accordionState.get(id)!;
    if (!this.accordionState.has(id)) return false;
    return false;
  }
}
