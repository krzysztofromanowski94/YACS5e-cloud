#!/usr/bin/python
# -*- coding: utf-8 -*-

import requests
import json
import sys
from pymongo import MongoClient

rest_url = 'http://dnd5eapi.co/api/'


def ability_scores(database):
    result = requests.get(rest_url + 'ability-scores').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']
        for skill_iterator in range(len(result_item['skills'])):
            del result_item['skills'][skill_iterator]['url']

        #print 'ability_scores', result_item

        database.ability_scores.update_one({'index': result_item['index'
                ]}, {'$set': result_item}, upsert=True)


def skills(database):
    result = requests.get(rest_url + 'skills').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']
        del result_item['ability_score']['url']

        #print 'skills', result_item

        database.skills.update_one({'index': result_item['index']},
                                   {'$set': result_item}, upsert=True)


def proficiencies(database):
    result = requests.get(rest_url + 'proficiencies').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        for classes_iterator in range(len(result_item['classes'])):
            result_item['classes'][classes_iterator]['index'] = \
                requests.get(result_item['classes'
                             ][classes_iterator]['url']).json()['index']
            del result_item['classes'][classes_iterator]['url']

        for races_iterator in range(len(result_item['races'])):
            result_item['races'][races_iterator]['index'] = \
                requests.get(result_item['races'][races_iterator]['url'
                             ]).json()['index']
            del result_item['races'][races_iterator]['url']

        #print 'proficiencies', result_item

        database.proficiencies.update_one({'index': result_item['index'
                ]}, {'$set': result_item}, upsert=True)


def languages(database):
    result = requests.get(rest_url + 'languages').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        #print 'languages', result_item

        database.languages.update_one({'index': result_item['index']},
                {'$set': result_item}, upsert=True)


def classes(database):
    result = requests.get(rest_url + 'classes').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        for proficiency_choices_iterator in \
            range(len(result_item['proficiency_choices'])):
            for proficiency_choices_from_iterator in \
                range(len(result_item['proficiency_choices'
                      ][proficiency_choices_iterator]['from'])):
                try:
                    result_item['proficiency_choices'
                                ][proficiency_choices_iterator]['from'
                            ][proficiency_choices_from_iterator]['index'
                            ] = \
                        requests.get(result_item['proficiency_choices'
                            ][proficiency_choices_iterator]['from'
                            ][proficiency_choices_from_iterator]['url'
                            ]).json()['index']
                    del result_item['proficiency_choices'
                                    ][proficiency_choices_iterator]['from'
                            ][proficiency_choices_from_iterator]['url']
                except KeyError:
                    pass

        # In Monk class there is quite big nesting. This is solving the issue, yet does't look fancy.

                try:
                    for proficiency_choices_from_from_iterator in \
                        range(len(result_item['proficiency_choices'
                              ][proficiency_choices_iterator]['from'
                              ][proficiency_choices_from_iterator]['from'
                              ])):
                        proficiency_choices_from_from_index = \
                            requests.get(result_item['proficiency_choices'
                                ][proficiency_choices_iterator]['from'
                                ][proficiency_choices_from_iterator]['from'
                                ][proficiency_choices_from_from_iterator]['url'
                                ]).json()
                        proficiency_choices_from_from_index = \
                            proficiency_choices_from_from_index['index']
                        result_item['proficiency_choices'
                                    ][proficiency_choices_iterator]['from'
                                ][proficiency_choices_from_iterator]['from'
                                ][proficiency_choices_from_from_iterator]['index'
                                ] = proficiency_choices_from_from_index
                        del result_item['proficiency_choices'
                                ][proficiency_choices_iterator]['from'
                                ][proficiency_choices_from_iterator]['from'
                                ][proficiency_choices_from_from_iterator]['url'
                                ]
                except KeyError:
                    pass

        for proficient_iterator in range(len(result_item['proficiencies'
                ])):
            del result_item['proficiencies'][proficient_iterator]['url']

        for saving_throws_iterator in \
            range(len(result_item['saving_throws'])):
            del result_item['saving_throws'
                            ][saving_throws_iterator]['url']

        class_levels = requests.get(result_item['class_levels']['url'
                                    ].lower()).json()
        for class_levels_iterator in range(len(class_levels)):
            del class_levels[class_levels_iterator]['_id']
            del class_levels[class_levels_iterator]['url']
            del class_levels[class_levels_iterator]['class']
            for class_levels_features_iterator in \
                range(len(class_levels[class_levels_iterator]['features'
                      ])):
                class_level_feature = \
                    requests.get(class_levels[class_levels_iterator]['features'
                                 ][class_levels_features_iterator]['url'
                                 ]).json()
                del class_levels[class_levels_iterator]['features'
                        ][class_levels_features_iterator]['url']
                class_levels[class_levels_iterator]['features'
                        ][class_levels_features_iterator]['index'] = \
                    class_level_feature['index']

            if class_levels[class_levels_iterator]['feature_choices'] \
                != []:
                for class_levels_feature_choices_iterator in \
                    range(len(class_levels[class_levels_iterator]['feature_choices'
                          ])):
                    class_level_feature_choices = \
                        requests.get(class_levels[class_levels_iterator]['feature_choices'
                            ][class_levels_feature_choices_iterator]['url'
                            ]).json()
                    del class_levels[class_levels_iterator]['feature_choices'
                            ][class_levels_feature_choices_iterator]['url'
                            ]
                    class_levels[class_levels_iterator]['feature_choices'
                            ][class_levels_feature_choices_iterator]['index'
                            ] = class_level_feature_choices['index']

        result_item['class_levels'] = class_levels

        for subclasses_iterator in range(len(result_item['subclasses'
                ])):
            subclass_index = requests.get(result_item['subclasses'
                    ][subclasses_iterator]['url']).json()['index']
            del result_item['subclasses'][subclasses_iterator]['url']
            result_item['subclasses'][subclasses_iterator]['index'] = \
                subclass_index

        try:
            result_item['spellcasting']['index'] = \
                requests.get(result_item['spellcasting']['url'
                             ]).json()['index']
            del result_item['spellcasting']['url']
        except KeyError:
            pass

        result_item['starting_equipment']['index'] = \
            requests.get(result_item['starting_equipment']['url'
                         ]).json()['index']
        del result_item['starting_equipment']['url']

        #print 'classes', result_item

        database.classes.update_one({'index': result_item['index']},
                                    {'$set': result_item}, upsert=True)


def subclasses(database):
    result = requests.get(rest_url + 'subclasses').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        result_item['class']['index'] = requests.get(result_item['class'
                ]['url']).json()['index']
        del result_item['class']['url']

        for features_iterator in range(len(result_item['features'])):
            result_item['features'][features_iterator]['index'] = \
                requests.get(result_item['features'
                             ][features_iterator]['url']).json()['index'
                    ]
            del result_item['features'][features_iterator]['url']

        try:
            for spells_iterator in range(len(result_item['spells'])):
                result_item['spells'][spells_iterator]['spell']['index'
                        ] = requests.get(result_item['spells'
                        ][spells_iterator]['spell']['url'
                        ]).json()['index']
                del result_item['spells'][spells_iterator]['spell'
                        ]['url']

                if result_item['spells'
                               ][spells_iterator]['spell_acquisition_method'
                        ]['url'] != '':
                    result_item['spells'
                                ][spells_iterator]['spell_acquisition_method'
                            ]['index'] = \
                        requests.get(result_item['spells'
                            ][spells_iterator]['spell_acquisition_method'
                            ]['url']).json()['index']
                    pass
                del result_item['spells'
                                ][spells_iterator]['spell_acquisition_method'
                        ]['url']

                for spells_prerequisites_iterator in \
                    range(len(result_item['spells'
                          ][spells_iterator]['prerequisites'])):
                    result_item['spells'
                                ][spells_iterator]['prerequisites'
                            ][spells_prerequisites_iterator]['index'] = \
                        requests.get(result_item['spells'
                            ][spells_iterator]['prerequisites'
                            ][spells_prerequisites_iterator]['url'
                            ]).json()['index']
                    del result_item['spells'
                                    ][spells_iterator]['prerequisites'
                            ][spells_prerequisites_iterator]['url']
        except KeyError:
            pass

        #print 'subclasses', result_item

        database.subclasses.update_one({'index': result_item['index']},
                {'$set': result_item}, upsert=True)


def features(database):
    result = requests.get(rest_url + 'features').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        result_item['class']['index'] = requests.get(result_item['class'
                ]['url']).json()['index']
        del result_item['class']['url']

        try:
            result_item['subclass']['index'] = \
                requests.get(result_item['subclass']['url'
                             ]).json()['index']
            del result_item['subclass']['url']
        except KeyError:
            pass

        try:
            for choice_from_iterator in range(len(result_item['choice'
                    ]['from'])):
                result_item['choice']['from'
                        ][choice_from_iterator]['index'] = \
                    requests.get(result_item['choice']['from'
                                 ][choice_from_iterator]['url'
                                 ]).json()['index']
                del result_item['choice']['from'
                        ][choice_from_iterator]['url']
        except KeyError:
            pass

        #print 'features', result_item

        database.features.update_one({'index': result_item['index']},
                {'$set': result_item}, upsert=True)


def races(database):
    result = requests.get(rest_url + 'races').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        for starting_proficiencies_iterator in \
            range(len(result_item['starting_proficiencies'])):

            if 'url' in result_item['starting_proficiencies'
                                    ][starting_proficiencies_iterator]:
                result_item['starting_proficiencies'
                            ][starting_proficiencies_iterator]['index'
                        ] = \
                    requests.get(result_item['starting_proficiencies'
                                 ][starting_proficiencies_iterator]['url'
                                 ]).json()['index']
                del result_item['starting_proficiencies'
                                ][starting_proficiencies_iterator]['url'
                        ]
            elif 'from' in result_item['starting_proficiencies'
                    ][starting_proficiencies_iterator]:

                for starting_proficiencies_from_iterator in \
                    range(len(result_item['starting_proficiencies'
                          ][starting_proficiencies_iterator]['from'])):
                    result_item['starting_proficiencies'
                                ][starting_proficiencies_iterator]['from'
                            ][starting_proficiencies_from_iterator]['index'
                            ] = \
                        requests.get(result_item['starting_proficiencies'
                            ][starting_proficiencies_iterator]['from'
                            ][starting_proficiencies_from_iterator]['url'
                            ]).json()['index']
                    del result_item['starting_proficiencies'
                                    ][starting_proficiencies_iterator]['from'
                            ][starting_proficiencies_from_iterator]['url'
                            ]

        for languages_iterator in range(len(result_item['languages'])):
            if 'url' in result_item['languages'][languages_iterator]:
                result_item['languages'][languages_iterator]['index'] = \
                    requests.get(result_item['languages'
                                 ][languages_iterator]['url'
                                 ]).json()['index']
                del result_item['languages'][languages_iterator]['url']
            elif 'from' in result_item['languages'][languages_iterator]:
                for languages_from_iterator in \
                    range(len(result_item['languages'
                          ][languages_iterator]['from'])):
                    result_item['languages'][languages_iterator]['from'
                            ][languages_from_iterator]['index'] = \
                        requests.get(result_item['languages'
                            ][languages_iterator]['from'
                            ][languages_from_iterator]['url'
                            ]).json()['index']
                    del result_item['languages'
                                    ][languages_iterator]['from'
                            ][languages_from_iterator]['url']

        try:
            for traits_iterator in range(len(result_item['traits'])):
                result_item['traits'][traits_iterator]['index'] = \
                    requests.get(result_item['traits'
                                 ][traits_iterator]['url'
                                 ]).json()['index']
                del result_item['traits'][traits_iterator]['url']
        except KeyError:
            pass

        try:
            for subraces_iterator in range(len(result_item['subraces'
                    ])):
                result_item['subraces'][subraces_iterator]['index'] = \
                    requests.get(result_item['subraces'
                                 ][subraces_iterator]['url'
                                 ]).json()['index']
                del result_item['subraces'][subraces_iterator]['url']
        except KeyError:
            pass

        #print 'races', result_item

        database.races.update_one({'index': result_item['index']},
                                  {'$set': result_item}, upsert=True)


def subraces(database):
    result = requests.get(rest_url + 'subraces').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        result_item['race']['index'] = requests.get(result_item['race'
                ]['url']).json()['index']
        del result_item['race']['url']

        for racial_traits_iterator in \
            range(len(result_item['racial_traits'])):
            if 'url' in result_item['racial_traits'
                                    ][racial_traits_iterator]:
                result_item['racial_traits'
                            ][racial_traits_iterator]['index'] = \
                    requests.get(result_item['racial_traits'
                                 ][racial_traits_iterator]['url'
                                 ]).json()['index']
                del result_item['racial_traits'
                                ][racial_traits_iterator]['url']
            elif 'from' in result_item['racial_traits'
                    ][racial_traits_iterator]:
                for racial_traits_from_iterator in \
                    range(len(result_item['racial_traits'
                          ][racial_traits_iterator]['from'])):
                    result_item['racial_traits'
                                ][racial_traits_iterator]['from'
                            ][racial_traits_from_iterator]['index'] = \
                        requests.get(result_item['racial_traits'
                            ][racial_traits_iterator]['from'
                            ][racial_traits_from_iterator]['url'
                            ]).json()['index']
                    del result_item['racial_traits'
                                    ][racial_traits_iterator]['from'
                            ][racial_traits_from_iterator]['url']

        try:
            try:
                result_item['starting_proficiencies'] = \
                    result_item['starting_proficiencies:']
                del result_item['starting_proficiencies:']
            except KeyError:
                pass
            for starting_proficiencies_iterator in \
                range(len(result_item['starting_proficiencies'])):
                result_item['starting_proficiencies'
                            ][starting_proficiencies_iterator]['index'
                        ] = \
                    requests.get(result_item['starting_proficiencies'
                                 ][starting_proficiencies_iterator]['url'
                                 ]).json()['index']
                del result_item['starting_proficiencies'
                                ][starting_proficiencies_iterator]['url'
                        ]
        except KeyError:
            pass

        try:
            for languages_iterator in range(len(result_item['languages'
                    ])):
                for languages_from_iterator in \
                    range(len(result_item['languages'
                          ][languages_iterator]['from'])):
                    result_item['languages'][languages_iterator]['from'
                            ][languages_from_iterator]['index'] = \
                        requests.get(result_item['languages'
                            ][languages_iterator]['from'
                            ][languages_from_iterator]['url'
                            ]).json()['index']
                    del result_item['languages'
                                    ][languages_iterator]['from'
                            ][languages_from_iterator]['url']
        except KeyError:
            pass

        #print 'subraces', result_item

        database.subraces.update_one({'index': result_item['index']},
                {'$set': result_item}, upsert=True)


def equipment(database):
    result = requests.get(rest_url + 'equipment').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        try:

# because in some cases semantic is not consistent

            if 'damage_type' in result_item['damage']:
                result_item['damage']['type'] = result_item['damage'
                        ]['damage_type']
                del result_item['damage']['damage_type']

# this is a workaround for error in restapi

            if type(result_item['damage']['type']['name']) is dict:
                result_item['damage']['type'] = result_item['damage'
                        ]['type']['name']
                del result_item['damage']['type']['name']
            result_item['damage']['type']['index'] = \
                requests.get(result_item['damage']['type']['url'
                             ]).json()['index']
            del result_item['damage']['type']['url']
        except KeyError:
            pass

        try:
            if 'damage_type' in result_item['2h_damage']:
                result_item['2h_damage']['type'] = \
                    result_item['2h_damage']['damage_type']
                del result_item['2h_damage']['damage_type']
            result_item['2h_damage']['type']['index'] = \
                requests.get(result_item['2h_damage']['type']['url'
                             ]).json()['index']
            del result_item['2h_damage']['type']['url']
        except KeyError:
            pass

        try:
            for properties_iterator in \
                range(len(result_item['properties'])):
                if type(result_item['properties'][properties_iterator]) \
                    is dict:
                    result_item['properties'
                                ][properties_iterator]['index'] = \
                        requests.get(result_item['properties'
                            ][properties_iterator]['url'
                            ]).json()['index']
                    del result_item['properties'
                                    ][properties_iterator]['url']
        except KeyError:
            pass

        try:
            if 'armor_category:' in result_item:
                result_item['armor_category'] = \
                    result_item['armor_category:']
                del result_item['armor_category:']
            if type(result_item['armor_category']) is dict:
                result_item['armor_category']['index'] = \
                    requests.get(result_item['armor_category']['url'
                                 ]).json()['index']
                del result_item['armor_category']['url']
        except KeyError:
            pass

        try:
            for contents_iterator in range(len(result_item['contents'
                    ])):
                result_item['contents'][contents_iterator]['name'] = \
                    requests.get(result_item['contents'
                                 ][contents_iterator]['item_url'
                                 ]).json()['name']
                result_item['contents'][contents_iterator]['index'] = \
                    requests.get(result_item['contents'
                                 ][contents_iterator]['item_url'
                                 ]).json()['index']
                del result_item['contents'
                                ][contents_iterator]['item_url']
        except KeyError:
            pass

        #print 'equipment', result_item

        database.equipment.update_one({'index': result_item['index']},
                {'$set': result_item}, upsert=True)


def conditions(database):
    result = requests.get(rest_url + 'conditions').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        #print 'conditions', result_item

        database.conditions.update_one({'index': result_item['index']},
                {'$set': result_item}, upsert=True)


def damage_types(database):
    result = requests.get(rest_url + 'damage-types').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        #print 'damage_types', result_item

        database.damage_types.update_one({'index': result_item['index'
                ]}, {'$set': result_item}, upsert=True)


def magic_schools(database):
    result = requests.get(rest_url + 'magic-schools').json()

    for iterator in range(result['count']):
        result_item = requests.get(result['results'][iterator]['url'
                                   ]).json()
        del result_item['_id']
        del result_item['url']

        #print 'magic_schools', result_item

        database.magic_schools.update_one({'index': result_item['index'
                ]}, {'$set': result_item}, upsert=True)


def main():
    database = MongoClient('localhost', 6969).mech
    ability_scores(database)
    skills(database)
    proficiencies(database)
    languages(database)
    classes(database)
    subclasses(database)
    features(database)
    races(database)
    subraces(database)
    equipment(database)
    conditions(database)
    damage_types(database)
    magic_schools(database)


if __name__ == '__main__':
    main()
