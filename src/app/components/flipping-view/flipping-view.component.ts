import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { Observable, Subscription } from 'rxjs';
import { CommonModule } from '@angular/common';
import { CdkAccordionModule } from '@angular/cdk/accordion';

@Component({
  selector: 'app-flipping-view',
  standalone: true,
  imports: [CommonModule, CdkAccordionModule],
  templateUrl: './flipping-view.component.html',
  styleUrl: './flipping-view.component.scss',
})
export class FlippingViewComponent implements OnInit, OnDestroy {
  @Input() listingData$!: Observable<IItemListings[]>;
  itemData: IItemListings[] = [];
  expandedItem: number = -1;
  private _subscription: Subscription = new Subscription();

  ngOnInit(): void {
    this._subscription = this.listingData$.subscribe((data) => {
      this.itemData = data.filter(
        (a) => a.margin > 0.25 && a.lowpricevolume + a.highpricevolume >= 25
      );
    });
  }

  ngOnDestroy(): void {
    this._subscription.unsubscribe();
  }

  isExpanded(id: number): boolean {
    return this.expandedItem === id;
  }

  expandItem(id: number): void {
    this.expandedItem = this.expandedItem === id ? -1 : id;
  }
}
