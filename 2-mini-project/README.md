# CPLT - Channel-based Concurrency Module


## Lab Class #2 - Mini-Project

To submit your answers, simply push your files onto the repository. The problem will be graded.

----

## Trains and Railway Tracks

When designing a railway system, careful planning is required due to the topographical constraints of the landscape 
where the railway tracks will be placed. In some locations, space allows for the placement of parallel railway tracks (i.e. double tracks), enabling 
trains to travel in opposite directions without using the same track. In such double track railways, the direction of each track is *fixed* and respected by all trains.
However, it is often the case that in certain terrains, only a single track can be placed, which must then be shared by trains 
going in both directions. 

A problem arises when a double track merges onto a single track. If a train is travelling East on the 
single track and another train is travelling West on the double track, the train on the double track will have to wait until
the Eastbound train passes the railway junction before it can enter the single track. 
The more general scenario is a double track junctioning into a single track which then connects to another double track, 
where many trains may need to be coordinated to ensure a safe use of the railway. 

In this Mini-Project, you will model three variants of these coordination mechanisms in Go, 
where the double track - single track - double track system is such that each of the tracks in the double track portions has a fixed direction 
that is respected by all trains and the single track can be used in both directions (but certainly not simultaneously!). 
This setup is fixed in the three tasks, detailed below.

### Task 1

In this task you will implement a simplified variant of the double track, single track, double track coordination scenario, where access to the single track (on both sides) is managed by a traffic light on each access point. The traffic light's control system is such that it automatically and instantaneously  synchronizes the lights on the two sides of the single track (as if it were a single semaphore system).

Your goal in this task is to implement a concurrent system with two trains (which run the same code), initially running on opposite directions but in different double tracks, and the single track control system. In this task, you can abstract the reality that the control system must internally synchronize the lights on the two sides of the single track, meaning that you can assume a single running entity (without any additional internal concurrency control) that manages access to both sides of the single track as if it were a single physical semaphore (e.g., as if the trains can see the lights on both sides).

Your system should have:
  * Two trains, running the same code, which you can may assume are running in opposite directions
  * The single-track control system, with which both trains interact.

The control system must never allow more than one train (in opposite directions) on the single track at once. The system has a whole must have **at least** two goroutines (excluding the main goroutine). Your control system should further be designed under the assumption that there are never more than two trains in opposing directions and that there is only one train in a given direction at a time.

*Note:* This task is meant to be straightforward. I am not trying to trick you.

You should not use any features provided by the Go ``sync`` package, with the exception of the use of ``WaitGroup`` described in lecture.

1. To submit your work, define a ``task1`` package (inside the ``task1`` folder in the repository) that exports a ``Setup()`` function that assembles and runs the system composed of the two trains and the single-track control system. You may want to define a ``Train`` function for this as well, where both trains can interact directly with the global control system.
2. Think critically of your work. Add a small report to the ``task1`` folder that explains how your system ensures that the control system never lets two trains enter in a collision course, given the assumptions above.


### Task 2

Now that you have designed a single-track control system under several simplifying assumptions, you are now ready to lift some of the simplifications and make it a bit more sophisticated.

In this Task, the implementation of the control system may no longer assume that the lights on both sides of the track are instantaneously synchronized and must implement an internal synchronization protocol for the lights at the two ends of the single track which guarantees it is never the case that two trains in opposite directions enter the single track.
Specifically, when a train arrives at a junction:
  - The train requests permission to enter the single track (effectively checking if the light is green);
  - If a train is already on the single track, the train must wait for the other to pass through.
  - If no train is on the track, the train can enter and the light on the opposite side stops trains from entering.
  - When the train passes through, it unlocks the ability for other trains to enter.
  - If two trains arrive "at the same time", only one can pass.

You may still work under the assumption that there is only one train traveling in a given direction at a time. 



1. To submit your work, define a ``task2`` package (inside the ``task2`` folder in the repository) that exports a ``Setup()`` function that assembles and runs the system composed of the two trains and the (now more sophisticated) single-track control system. Your system now should have the each train interacting with a separate light (which allows or disallows the train from moving), which then interact with the system ensure safety of the trains. You may have each train notify their respective "light" that they have cleared the single track, for the sake of simplicity.
2. Add a small report to the ``task2`` folder that explains what correctness properties your system has and what achieves them.

### Task 3

We now lift the assumption that only one train travels in a given direction at a time. Redesign the control system such that it now is aware of the direction of the train(s) currently using it (if any). When a new train arrives at one of the two junctions, if another train is traveling in the single track **in the same direction** as the new train, it may enter the single track. You may assume the single track can fit an arbitrary number of trains. If any number of trains are trave  ling in the single track in the opposite direction, the train must wait for all of them. You may also assume that the double track can also have an unbounded number of trains waiting to enter the single-track, but they should enter the single-track *in the order they arrive*.

1. To submit your work, define a ``task3`` package (inside the ``task3`` folder in the repository) that exports a ``Setup()`` function that assembles and runs the system composed of 6 trains (3 in each direction) and the (now more sophisticated) single-track control system. The system should be an extension of the setup of Task 2.

2. Add a small report to the ``task3`` folder that explains what correctness properties your system has and what achieves them.

---

## General Advice

- Do not use locks (i.e. those in the ``sync`` package)!
- It will be useful if the various components (at least the trains) log their actions. For instance, the train logs when it contacts the system/light to enter the single-track, when it actually enters and when it leaves.
- Keep things as simple as possible. Some abstractions may require only sharing a few (or even one) channel(s), others require a handler goroutine that "multiplexes" requests from various channels.