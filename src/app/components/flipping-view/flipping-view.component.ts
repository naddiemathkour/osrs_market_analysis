import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { Observable, Subscription } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-flipping-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './flipping-view.component.html',
  styleUrl: './flipping-view.component.scss',
})
export class FlippingViewComponent implements OnInit, OnDestroy {
  @Input() listingData$!: Observable<IItemListings[]>;
  itemData: IItemListings[] = [];
  private _subscription: Subscription = new Subscription();

  ngOnInit(): void {
    this._subscription = this.listingData$.subscribe(
      (data) => (this.itemData = data.filter((a) => a.margin > 0.1))
    );
  }

  ngOnDestroy(): void {
    this._subscription.unsubscribe();
  }
}
