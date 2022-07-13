import { useState } from "react";
import { useIntl } from "react-intl";
import { AxiosResponse } from "axios";
import {
  closestCenter,
  DndContext,
  DragOverlay,
  useSensor,
  useSensors,
  PointerSensor,
  KeyboardSensor,
  useDndContext,
  MeasuringStrategy,
  DropAnimation,
  defaultDropAnimationSideEffects,
} from "@dnd-kit/core";
import type {
  Active,
  Over,
  DragStartEvent,
  DragEndEvent,
  MeasuringConfiguration,
} from "@dnd-kit/core";
import { KeyedMutator } from "swr";
import {
  arrayMove,
  useSortable,
  SortableContext,
  sortableKeyboardCoordinates,
} from "@dnd-kit/sortable";
import { CSS, isKeyboardEvent } from "@dnd-kit/utilities";

import { LaxoProductAsset, LaxoProductDetails } from "@/types/ApiResponse";
import useProductApi from "@/hooks/useProductApi";
import {
  AssetManagementItem,
  AssetManagementItemProps,
  Position,
} from "@/components/dashboard/product/AssetManagement/AssetManagementItem";
import { useDashboard } from "@/providers/DashboardProvider";

type LaxoActive = Omit<Active, "id"> & {
  id: string;
};

type LaxoOver = Omit<Over, "id"> & {
  id: string;
};

type LaxoDragStartEvent = Omit<DragStartEvent, "active"> & {
  active: LaxoActive;
};

type LaxoDragEndEvent = Omit<DragEndEvent, "over"> & {
  over: LaxoOver | null;
};

const dropAnimation: DropAnimation = {
  keyframes({ transform }) {
    return [
      { transform: CSS.Transform.toString(transform.initial) },
      {
        transform: CSS.Transform.toString({
          scaleX: 0.98,
          scaleY: 0.98,
          x: transform.final.x - 10,
          y: transform.final.y - 10,
        }),
      },
    ];
  },
  sideEffects: defaultDropAnimationSideEffects({
    styles: {
      active: { opacity: "0" },
    },
  }),
};

const measuring: MeasuringConfiguration = {
  droppable: {
    strategy: MeasuringStrategy.Always,
  },
};

type DragAndDropContainerProps = {
  assets: LaxoProductAsset[];
  assetsToken: string;
  setActiveAssetDetails: (arg: LaxoProductAsset) => void;
  setShowImageDetails: (arg: boolean) => void;
  productID: string;
  detailMutate: KeyedMutator<AxiosResponse<LaxoProductDetails, any>>;
};

export default function DragAndDropContainer({
  assets,
  assetsToken,
  setActiveAssetDetails,
  setShowImageDetails,
  productID,
  detailMutate,
}: DragAndDropContainerProps) {
  const { doChangeImageOrder } = useProductApi();
  const t = useIntl();

  const [activeId, setActiveId] = useState<string | null>(null);
  const [items, setItems] = useState<LaxoProductAsset[]>(assets);

  const { dashboardDispatch } = useDashboard();

  const activeIndex = activeId
    ? items
        .map(function (e) {
          return e.id;
        })
        .indexOf(activeId)
    : -1;

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  async function handleOrderChange(items: LaxoProductAsset[]) {
    const newOrder = items.map((item, index) => ({
      assetID: item.id,
      order: index,
    }));

    const result = await doChangeImageOrder(productID, newOrder);
    if (result) {
      detailMutate();
      dashboardDispatch({
        type: "alert",
        alert: {
          type: "success",
          message: t.formatMessage({
            description: "Change image order success",
            defaultMessage: "Successfully updated your image order",
          }),
        },
      });
    } else {
      dashboardDispatch({
        type: "alert",
        alert: {
          type: "error",
          message: t.formatMessage({
            description: "Change image order server error",
            defaultMessage:
              "Something went wrong while removing your image, try again later",
          }),
        },
      });
    }
  }

  function handleDragStart({ active }: LaxoDragStartEvent) {
    setActiveId(active.id);
  }

  function handleDragEnd({ over }: LaxoDragEndEvent) {
    if (over) {
      const overIndex = items
        .map(function (e) {
          return e.id;
        })
        .indexOf(over.id);

      if (activeIndex !== overIndex) {
        const newIndex = overIndex;

        setItems((items) => {
          const newArray = arrayMove(items, activeIndex, newIndex);
          handleOrderChange(newArray);
          return newArray;
        });
      }
    }

    setActiveId(null);
  }

  function handleDragCancel() {
    setActiveId(null);
  }

  return (
    <DndContext
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
      onDragCancel={handleDragCancel}
      sensors={sensors}
      collisionDetection={closestCenter}
      measuring={measuring}
    >
      <SortableContext items={items}>
        <ul role="list" className="grid list-none grid-cols-4 gap-x-4 gap-y-8">
          {items.map((m) => (
            <SortableItem
              activeIndex={activeIndex}
              key={m.id}
              asset={m}
              setShowImageDetails={setShowImageDetails}
              setActiveAssetDetails={setActiveAssetDetails}
              assetsToken={assetsToken}
            />
          ))}
        </ul>
      </SortableContext>
      <DragOverlay dropAnimation={dropAnimation}>
        {activeId ? (
          <ul className="list-none">
            <AssetManagementItem
              clone
              asset={items[activeIndex]}
              setShowImageDetails={setShowImageDetails}
              setActiveAssetDetails={setActiveAssetDetails}
              assetsToken={assetsToken}
            />
          </ul>
        ) : null}
      </DragOverlay>
    </DndContext>
  );
}

function SortableItem({
  activeIndex,
  ...props
}: AssetManagementItemProps & { activeIndex: number }) {
  const {
    attributes,
    listeners,
    index,
    isDragging,
    isSorting,
    over,
    setNodeRef,
    transform,
    transition,
  } = useSortable({
    id: props.asset.id,
    animateLayoutChanges: always,
  });

  return (
    <AssetManagementItem
      ref={setNodeRef}
      active={isDragging}
      style={{
        transition,
        transform: isSorting ? undefined : CSS.Translate.toString(transform),
      }}
      insertPosition={
        over?.id === props.asset.id
          ? index > activeIndex
            ? Position.After
            : Position.Before
          : undefined
      }
      {...props}
      {...attributes}
      {...listeners}
    />
  );
}

function always() {
  return true;
}
