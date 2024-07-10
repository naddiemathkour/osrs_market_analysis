import { Component, OnInit } from '@angular/core';
import { FetchmarketdataService } from '../../services/fetchmarketdata.service';
import { Observable } from 'rxjs';
import { IItemListings } from '../../interfaces/itemlistings.interface';
import { CommonModule } from '@angular/common';
import { FlippingViewComponent } from '../flipping-view/flipping-view.component';
import { AlchingViewComponent } from '../alching-view/alching-view.component';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    FlippingViewComponent,
    AlchingViewComponent,
  ],
  providers: [FetchmarketdataService],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent implements OnInit {
  listingData$: Observable<IItemListings[]> = new Observable<IItemListings[]>();
  test$!: Observable<IItemListings[]>;
  flippingview: boolean = false;
  alchingview: boolean = false;
  showMembers: boolean = true;

  constructor(private _fetchMarketData: FetchmarketdataService) {}

  ngOnInit() {
    // Create cold observable to be subscribed to in child components
    this.listingData$ = this._fetchMarketData.getMarketData();
  }

  toggleFlippingView(): void {
    this.alchingview = false;
    this.flippingview = true;
  }

  toggleAlchingView(): void {
    this.flippingview = false;
    this.alchingview = true;
  }
}
