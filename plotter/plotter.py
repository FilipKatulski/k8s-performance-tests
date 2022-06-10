import csv
import os, subprocess 
import argparse
from time import sleep
import matplotlib.pyplot as plt
from pip import main

try:
    from art import text2art
except ModuleNotFoundError:
    print("\n\n'Art' module is not installed\n")
    pass


def display_help():
    """
    Displaying custom help message
    """
    ...


def display_header():
    """
    Displays python art header
    """
    try:
        art_1 = text2art("ClusterLoader2 plotting module")
        print(art_1)
        sleep(1)
    except:
        print("\nClusterLoader2 plotting moduler\n\n")
        sleep(1)


def data_parsing():
    """
    Function responsible for timeline parsing
    """
    ...


def plot_selection():
    """
    Function to select which plot user wants to create
    """
    ...


def main():
    """
    The main function responsible for plotting 
    """
    ...


if __name__ == "main":
    main()
